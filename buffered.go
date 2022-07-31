// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memnet

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iangudger/memnet/unix"
	"github.com/iangudger/memnet/unix/linux"
	"github.com/iangudger/memnet/unix/syserr"
	"github.com/iangudger/memnet/unix/tcpip"
	"github.com/iangudger/memnet/unix/waiter"
)

var (
	errCanceled   = errors.New("operation canceled")
	errWouldBlock = errors.New("operation would block")
)

// timeoutError is how the net package reports timeouts.
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

type deadlineTimer struct {
	// mu protects the fields below.
	mu sync.Mutex

	readTimer     *time.Timer
	readCancelCh  chan struct{}
	writeTimer    *time.Timer
	writeCancelCh chan struct{}
}

func (d *deadlineTimer) init() {
	d.readCancelCh = make(chan struct{})
	d.writeCancelCh = make(chan struct{})
}

func (d *deadlineTimer) readCancel() <-chan struct{} {
	d.mu.Lock()
	c := d.readCancelCh
	d.mu.Unlock()
	return c
}
func (d *deadlineTimer) writeCancel() <-chan struct{} {
	d.mu.Lock()
	c := d.writeCancelCh
	d.mu.Unlock()
	return c
}

// setDeadline contains the shared logic for setting a deadline.
//
// cancelCh and timer must be pointers to deadlineTimer.readCancelCh and
// deadlineTimer.readTimer or deadlineTimer.writeCancelCh and
// deadlineTimer.writeTimer.
//
// setDeadline must only be called while holding d.mu.
func (d *deadlineTimer) setDeadline(cancelCh *chan struct{}, timer **time.Timer, t time.Time) {
	if *timer != nil && !(*timer).Stop() {
		*cancelCh = make(chan struct{})
	}

	// Create a new channel if we already closed it due to setting an already
	// expired time. We won't race with the timer because we already handled
	// that above.
	select {
	case <-*cancelCh:
		*cancelCh = make(chan struct{})
	default:
	}

	// "A zero value for t means I/O operations will not time out."
	// - net.Conn.SetDeadline
	if t.IsZero() {
		return
	}

	timeout := t.Sub(time.Now())
	if timeout <= 0 {
		close(*cancelCh)
		return
	}

	// Timer.Stop returns whether or not the AfterFunc has started, but
	// does not indicate whether or not it has completed. Make a copy of
	// the cancel channel to prevent this code from racing with the next
	// call of setDeadline replacing *cancelCh.
	ch := *cancelCh
	*timer = time.AfterFunc(timeout, func() {
		close(ch)
	})
}

// SetReadDeadline implements net.Conn.SetReadDeadline and
// net.PacketConn.SetReadDeadline.
func (d *deadlineTimer) SetReadDeadline(t time.Time) error {
	d.mu.Lock()
	d.setDeadline(&d.readCancelCh, &d.readTimer, t)
	d.mu.Unlock()
	return nil
}

// SetWriteDeadline implements net.Conn.SetWriteDeadline and
// net.PacketConn.SetWriteDeadline.
func (d *deadlineTimer) SetWriteDeadline(t time.Time) error {
	d.mu.Lock()
	d.setDeadline(&d.writeCancelCh, &d.writeTimer, t)
	d.mu.Unlock()
	return nil
}

// SetDeadline implements net.Conn.SetDeadline and net.PacketConn.SetDeadline.
func (d *deadlineTimer) SetDeadline(t time.Time) error {
	d.mu.Lock()
	d.setDeadline(&d.readCancelCh, &d.readTimer, t)
	d.setDeadline(&d.writeCancelCh, &d.writeTimer, t)
	d.mu.Unlock()
	return nil
}

// uidGen implements UniqueIDProvider.
type uidGen struct {
	// id must be accessed atomically.
	id uint64
}

// UniqueID implements UniqueIDProvider.UniqueID.
func (ug *uidGen) UniqueID() uint64 {
	return atomic.AddUint64(&ug.id, 1)
}

// uniqueIDProvider is the shared UniqueIDProvider.
var uniqueIDProvider uidGen

// A BufferedStreamConn is an in-memory Unix stream socket emulator that
// implements the net.Conn interface.
//
// Emulates the "unix" network in the net package (SOCK_STREAM), message
// boundaries are not preserved, but bytes are transferred reliably and in the
// order that they were sent. ReadXxx calls will read the next available bytes
// and any remaining unread bytes will be available to the next ReadXxx call.
type BufferedStreamConn struct {
	deadlineTimer

	ep unix.Endpoint
}

var _ net.Conn = (*BufferedStreamConn)(nil)

// NewBufferedStreamConnPair creates a connected pair (socketpair) of
// BufferedStreamConns.
func NewBufferedStreamConnPair() (*BufferedStreamConn, *BufferedStreamConn) {
	ep1, ep2 := unix.NewPair(nil, linux.SOCK_STREAM, &uniqueIDProvider)

	c1 := &BufferedStreamConn{
		ep: ep1,
	}
	c1.deadlineTimer.init()

	c2 := &BufferedStreamConn{
		ep: ep2,
	}
	c2.deadlineTimer.init()
	return c1, c2
}

type opErrorer interface {
	newOpError(op string, err error) *net.OpError
}

// commonRead implements the common logic between net.Conn.Read and
// net.PacketConn.ReadFrom.
func commonRead(ep unix.Endpoint, deadline <-chan struct{}, data [][]byte, addr *tcpip.FullAddress, errorer opErrorer, dontWait bool) (int64, error) {
	select {
	case <-deadline:
		return 0, errorer.newOpError("read", &timeoutError{})
	default:
	}

	read, _ /* msgLen */, _ /* cm */, _ /* CMTruncated */, err := ep.RecvMsg(nil /* peek */, data, false /* creds */, 0 /* numRights */, false /* peek */, addr)
	if err == syserr.ErrWouldBlock {
		if dontWait {
			return 0, errWouldBlock
		}
		// Create wait queue entry that notifies a channel.
		waitEntry, notifyCh := waiter.NewChannelEntry(nil)
		ep.EventRegister(&waitEntry, waiter.EventIn)
		defer ep.EventUnregister(&waitEntry)
		for {
			read, _ /* msgLen */, _ /* cm */, _ /* CMTruncated */, err = ep.RecvMsg(nil /* peek */, data, false /* creds */, 0 /* numRights */, false /* peek */, addr)
			if err != syserr.ErrWouldBlock {
				break
			}

			select {
			case <-deadline:
				return 0, errorer.newOpError("read", &timeoutError{})
			case <-notifyCh:
			}
		}
	}

	if err == syserr.ErrClosedForReceive {
		return 0, io.EOF
	}

	if err != nil {
		return 0, errorer.newOpError("read", errors.New(err.String()))
	}

	return read, nil
}

// Read implements net.Conn.Read.
func (c *BufferedStreamConn) Read(b []byte) (int, error) {
	deadline := c.readCancel()

	vec := [][]byte{b}
	var err error
	for err == nil && len(vec[0]) > 0 {
		_, err = commonRead(c.ep, deadline, vec, nil /* addr */, c, len(vec[0]) != len(b) /* dontWait */)
	}

	if err != nil && len(vec[0]) == len(b) {
		return 0, err
	}

	return int(len(b) - len(vec[0])), nil
}

// Write implements net.Conn.Write.
func (c *BufferedStreamConn) Write(b []byte) (int, error) {
	deadline := c.writeCancel()

	// Check if deadlineTimer has already expired.
	select {
	case <-deadline:
		return 0, c.newOpError("write", &timeoutError{})
	default:
	}

	vec := [][]byte{b}

	// We must handle two soft failure conditions simultaneously:
	//  1. Write may write nothing and return syserr.ErrWouldBlock.
	//     If this happens, we need to register for notifications if we have
	//     not already and wait to try again.
	//  2. Write may write fewer than the full number of bytes and return
	//     without error. In this case we need to try writing the remaining
	//     bytes again. I do not need to register for notifications.
	//
	// What is more, these two soft failure conditions can be interspersed.
	// There is no guarantee that all of the condition #1s will occur before
	// all of the condition #2s or visa-versa.
	var (
		err      *syserr.Error
		nbytes   int
		reg      bool
		notifyCh chan struct{}
	)
	for nbytes < len(b) && (err == syserr.ErrWouldBlock || err == nil) {
		if err == syserr.ErrWouldBlock {
			if !reg {
				// Only register once.
				reg = true

				// Create wait queue entry that notifies a channel.
				var waitEntry waiter.Entry
				waitEntry, notifyCh = waiter.NewChannelEntry(nil)
				c.ep.EventRegister(&waitEntry, waiter.EventOut)
				defer c.ep.EventUnregister(&waitEntry)
			} else {
				// Don't wait immediately after registration in case more data
				// became available between when we last checked and when we setup
				// the notification.
				select {
				case <-deadline:
					return nbytes, c.newOpError("write", &timeoutError{})
				case <-notifyCh:
				}
			}
		}

		var n int64
		n, err = c.ep.SendMsg(nil /* ctx */, vec, unix.ControlMessages{}, nil /* BoundEndpoint */)
		nbytes += int(n)
		vec[0] = vec[0][n:]
	}

	if err == nil {
		return nbytes, nil
	}

	return nbytes, c.newOpError("write", errors.New(err.String()))
}

// Close implements net.Conn.Close.
func (c *BufferedStreamConn) Close() error {
	c.ep.Close(nil /* ctx */)
	return nil
}

// CloseRead shuts down the reading side of the connection. Most callers
// should just use Close.
//
// A Half-Close is performed the same as CloseRead for *net.UnixConn.
func (c *BufferedStreamConn) CloseRead() error {
	if terr := c.ep.Shutdown(tcpip.ShutdownRead); terr != nil {
		return c.newOpError("close", errors.New(terr.String()))
	}
	return nil
}

// CloseWrite shuts down the writing side of the connection. Most callers
// should just use Close.
//
// A Half-Close is performed the same as CloseWrite for *net.UnixConn.
func (c *BufferedStreamConn) CloseWrite() error {
	if terr := c.ep.Shutdown(tcpip.ShutdownWrite); terr != nil {
		return c.newOpError("close", errors.New(terr.String()))
	}
	return nil
}

// LocalAddr implements net.Conn.LocalAddr.
func (c *BufferedStreamConn) LocalAddr() net.Addr {
	a, err := c.ep.GetLocalAddress()
	if err != nil {
		return nil
	}
	return &net.UnixAddr{
		Net:  "unix",
		Name: string(a.Addr),
	}
}

// RemoteAddr implements net.Conn.RemoteAddr.
func (c *BufferedStreamConn) RemoteAddr() net.Addr {
	a, err := c.ep.GetRemoteAddress()
	if err != nil {
		return nil
	}
	return &net.UnixAddr{
		Net:  "unix",
		Name: string(a.Addr),
	}
}

func (c *BufferedStreamConn) newOpError(op string, err error) *net.OpError {
	return &net.OpError{
		Op:     op,
		Net:    "unix",
		Source: c.LocalAddr(),
		Addr:   c.RemoteAddr(),
		Err:    err,
	}
}

// A BufferedPacketConn is an in-memory Unix dgram/seqpacket socket emulator
// that implements the net.Conn interface.
//
// Effectively equivalent to the "unixgram" and "unixpacket" networks (SOCK_DGRAM and SOCK_SEQPACKET socket types respectively), message
// boundaries are not preserved, but bytes are transferred reliably and in the
// order that they were sent. ReadXxx calls will read the next available bytes
// and any remaining unread bytes will be available to the next ReadXxx call.
type BufferedPacketConn struct {
	deadlineTimer

	network string
	ep      unix.Endpoint
}

var (
	_ net.Conn       = (*BufferedPacketConn)(nil)
	_ net.PacketConn = (*BufferedPacketConn)(nil)
)

// NewBufferedPacketConnPair creates a connected pair (socketpair) of
// BufferedPacketConns.
//
// Supported networks: "unixgram", "unixpacket".
func NewBufferedPacketConnPair(network string) (*BufferedPacketConn, *BufferedPacketConn, error) {
	var stype linux.SockType
	switch network {
	case "unixgram":
		stype = linux.SOCK_DGRAM
	case "unixpacket":
		stype = linux.SOCK_SEQPACKET
	default:
		return nil, nil, &net.OpError{Op: "dial", Net: network, Source: &net.UnixAddr{Net: network}, Addr: &net.UnixAddr{Net: network}, Err: net.UnknownNetworkError(network)}
	}
	ep1, ep2 := unix.NewPair(nil, stype, &uniqueIDProvider)

	c1 := &BufferedPacketConn{
		network: network,
		ep:      ep1,
	}
	c1.deadlineTimer.init()

	c2 := &BufferedPacketConn{
		network: network,
		ep:      ep2,
	}
	c2.deadlineTimer.init()
	return c1, c2, nil
}

// LocalAddr implements net.Conn.LocalAddr.
func (c *BufferedPacketConn) LocalAddr() net.Addr {
	a, err := c.ep.GetLocalAddress()
	if err != nil {
		return nil
	}
	return &net.UnixAddr{
		Net:  c.network,
		Name: string(a.Addr),
	}
}

// RemoteAddr implements net.Conn.RemoteAddr.
func (c *BufferedPacketConn) RemoteAddr() net.Addr {
	a, err := c.ep.GetRemoteAddress()
	if err != nil {
		return nil
	}
	return &net.UnixAddr{
		Net:  c.network,
		Name: string(a.Addr),
	}
}

func (c *BufferedPacketConn) newOpError(op string, err error) *net.OpError {
	return c.newRemoteOpError(op, nil, err)
}

func (c *BufferedPacketConn) newRemoteOpError(op string, remote net.Addr, err error) *net.OpError {
	if remote == nil {
		remote = c.RemoteAddr()
	}
	return &net.OpError{
		Op:     op,
		Net:    c.network,
		Source: c.LocalAddr(),
		Addr:   remote,
		Err:    err,
	}
}

// Read implements net.Conn.Read
func (c *BufferedPacketConn) Read(b []byte) (int, error) {
	bytesRead, _, err := c.ReadFrom(b)
	return bytesRead, err
}

// ReadFrom implements net.PacketConn.ReadFrom.
func (c *BufferedPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	deadline := c.readCancel()

	vec := [][]byte{b}
	var addr tcpip.FullAddress
	read, err := commonRead(c.ep, deadline, vec, &addr, c, false)
	if err != nil {
		return 0, nil, err
	}

	return int(read), &net.UnixAddr{
		Net:  c.network,
		Name: string(addr.Addr),
	}, nil
}

func (c *BufferedPacketConn) Write(b []byte) (int, error) {
	deadline := c.writeCancel()

	// Check if deadline has already expired.
	select {
	case <-deadline:
		return 0, c.newOpError("write", &timeoutError{})
	default:
	}

	vec := [][]byte{b}
	n, err := c.ep.SendMsg(nil /* ctx */, vec, unix.ControlMessages{}, nil /* BoundEndpoint */)

	if err == syserr.ErrWouldBlock {
		// Create wait queue entry that notifies a channel.
		waitEntry, notifyCh := waiter.NewChannelEntry(nil)
		c.ep.EventRegister(&waitEntry, waiter.EventOut)
		defer c.ep.EventUnregister(&waitEntry)
		for {
			select {
			case <-deadline:
				return int(n), c.newOpError("write", &timeoutError{})
			case <-notifyCh:
			}

			n, err = c.ep.SendMsg(nil /* ctx */, vec, unix.ControlMessages{}, nil /* BoundEndpoint */)
			if err != syserr.ErrWouldBlock {
				break
			}
		}
	}

	if err == nil {
		return int(n), nil
	}

	return int(n), c.newOpError("write", errors.New(err.String()))
}

// WriteTo implements net.PacketConn.WriteTo.
func (c *BufferedPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	if addr != nil {
		// TODO: Support address namespace.
		return 0, c.newRemoteOpError("write", addr, errors.New(syserr.ErrConnectionRefused.String()))
	}

	return c.Write(b)
}

// Close implements net.PacketConn.Close.
func (c *BufferedPacketConn) Close() error {
	c.ep.Close(nil /* ctx */)
	return nil
}

// CloseRead shuts down the reading side of the connection. Most callers
// should just use Close.
//
// A Half-Close is performed the same as CloseRead for *net.UnixConn.
func (c *BufferedPacketConn) CloseRead() error {
	if terr := c.ep.Shutdown(tcpip.ShutdownRead); terr != nil {
		return c.newOpError("close", errors.New(terr.String()))
	}
	return nil
}

// CloseWrite shuts down the writing side of the connection. Most callers
// should just use Close.
//
// A Half-Close is performed the same as CloseWrite for *net.UnixConn.
func (c *BufferedPacketConn) CloseWrite() error {
	if terr := c.ep.Shutdown(tcpip.ShutdownWrite); terr != nil {
		return c.newOpError("close", errors.New(terr.String()))
	}
	return nil
}

// NewBufferedPair creates a connected pair (socketpair) of
// buffered in-memory net.Conns.
//
// Supported networks: "unix", "unixgram", "unixpacket".
func NewBufferedPair(network string) (net.Conn, net.Conn, error) {
	switch network {
	case "unix":
		c1, c2 := NewBufferedStreamConnPair()
		return c1, c2, nil
	case "unixgram", "unixpacket":
		return NewBufferedPacketConnPair(network)
	default:
		return nil, nil, &net.OpError{Op: "dial", Net: network, Source: &net.UnixAddr{Net: network}, Addr: &net.UnixAddr{Net: network}, Err: net.UnknownNetworkError(network)}
	}
}

// A BufferedListener is an in-memory Unix socket emulator that implements
// the net.Listener interface.
type BufferedListener struct {
	ep         unix.Endpoint
	network    string
	cancelOnce sync.Once
	cancel     chan struct{}
}

var _ net.Listener = (*BufferedListener)(nil)

// BufferedListen creates a new BufferedListener not in an address namespace.
//
// Supported networks: "unix", "unixpacket".
func BufferedListen(network string, addr *net.UnixAddr) (*BufferedListener, error) {
	var ep unix.Endpoint

	if addr == nil {
		return nil, &net.OpError{Op: "listen", Net: network, Err: net.InvalidAddrError("nil")}
	}

	switch network {
	case "unix":
		ep = unix.NewConnectioned(nil /* ctx */, linux.SOCK_STREAM, &uniqueIDProvider)
	case "unixpacket":
		ep = unix.NewConnectioned(nil /* ctx */, linux.SOCK_SEQPACKET, &uniqueIDProvider)
	default:
		return nil, &net.OpError{Op: "listen", Net: network, Source: addr, Err: net.UnknownNetworkError(network)}
	}

	if err := ep.Bind(tcpip.FullAddress{Addr: tcpip.Address(addr.Name)}, nil /* commit */); err != nil {
		ep.Close(nil /* ctx */)
		return nil, &net.OpError{Op: "bind", Net: network, Source: addr, Err: errors.New(err.String())}
	}

	if err := ep.Listen(10 /* backlog */); err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: addr, Err: errors.New(err.String())}
	}

	return &BufferedListener{ep: ep, network: network, cancel: make(chan struct{})}, nil
}

// Close implements net.Listener.Close.
func (l *BufferedListener) Close() error {
	l.Shutdown()
	l.ep.Close(nil /* ctx */)
	return nil
}

// Shutdown stops the listener.
func (l *BufferedListener) Shutdown() {
	l.ep.Shutdown(tcpip.ShutdownWrite | tcpip.ShutdownRead)
	l.cancelOnce.Do(func() {
		close(l.cancel) // broadcast cancellation
	})
}

// Addr implements net.Listener.Addr.
func (l *BufferedListener) Addr() net.Addr {
	a, err := l.ep.GetLocalAddress()
	if err != nil {
		return nil
	}
	return &net.UnixAddr{
		Net:  l.network,
		Name: string(a.Addr),
	}
}

// Accept implements net.Conn.Accept.
func (l *BufferedListener) Accept() (net.Conn, error) {
	ep, err := l.ep.Accept()

	if err == syserr.ErrWouldBlock {
		// Create wait queue entry that notifies a channel.
		waitEntry, notifyCh := waiter.NewChannelEntry(nil)
		l.ep.EventRegister(&waitEntry, waiter.EventIn)
		defer l.ep.EventUnregister(&waitEntry)

		for {
			ep, err = l.ep.Accept()

			if err != syserr.ErrWouldBlock {
				break
			}

			select {
			case <-l.cancel:
				return nil, errCanceled
			case <-notifyCh:
			}
		}
	}

	if err != nil {
		return nil, &net.OpError{
			Op:   "accept",
			Net:  l.network,
			Addr: l.Addr(),
			Err:  errors.New(err.String()),
		}
	}

	if l.network == "unix" {
		c := &BufferedStreamConn{
			ep: ep,
		}
		c.deadlineTimer.init()
		return c, nil
	}

	c := &BufferedPacketConn{
		network: l.network,
		ep:      ep,
	}
	c.deadlineTimer.init()
	return c, nil
}

// Dial creates a new net.Conn connected to the listener and optionally bound
// to laddr.
func (l *BufferedListener) Dial(laddr *net.UnixAddr) (net.Conn, error) {
	return l.DialContext(context.Background(), laddr)
}

// DialContext creates a new net.Conn connected to the listener and optionally
// bound to laddr with the option of adding cancellation and timeouts.
func (l *BufferedListener) DialContext(ctx context.Context, laddr *net.UnixAddr) (net.Conn, error) {
	var ep unix.Endpoint

	switch l.network {
	case "unix":
		ep = unix.NewConnectioned(nil /* ctx */, linux.SOCK_STREAM, &uniqueIDProvider)
	case "unixpacket":
		ep = unix.NewConnectioned(nil /* ctx */, linux.SOCK_SEQPACKET, &uniqueIDProvider)
	default:
		return nil, &net.OpError{Op: "dial", Net: l.network, Source: laddr, Addr: l.Addr(), Err: net.UnknownNetworkError(l.network)}
	}

	bep, ok := l.ep.(unix.BoundEndpoint)
	if !ok {
		return nil, &net.OpError{Op: "dial", Net: l.network, Source: laddr, Addr: l.Addr(), Err: errors.New(syserr.ErrInvalidArgument.String())}
	}

	var (
		reg      bool
		notifyCh chan struct{}
	)

	for {
		err := ep.Connect(nil /* ctx */, bep)
		if err == nil {
			if l.network == "unix" {
				c := &BufferedStreamConn{
					ep: ep,
				}
				c.deadlineTimer.init()
				return c, nil
			}

			c := &BufferedPacketConn{
				network: l.network,
				ep:      ep,
			}
			c.deadlineTimer.init()
			return c, nil
		}
		if err != syserr.ErrWouldBlock {
			return nil, &net.OpError{Op: "dial", Net: l.network, Source: laddr, Addr: l.Addr(), Err: errors.New(syserr.ErrInvalidArgument.String())}
		}

		if !reg {
			// Only register once.
			reg = true

			// Create wait queue entry that notifies a channel.
			var waitEntry waiter.Entry
			waitEntry, notifyCh = waiter.NewChannelEntry(nil)
			l.ep.EventRegister(&waitEntry, waiter.EventOut)
			defer l.ep.EventUnregister(&waitEntry)

			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-notifyCh:
		}
	}
}
