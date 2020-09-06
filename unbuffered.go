// Copyright 2020 Ian Gudger
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
	"net"
)

type unbufferedPacketConn struct {
	net.Conn
}

var _ net.PacketConn = (*unbufferedPacketConn)(nil)

// ReadFrom implements net.PacketConn.ReadFrom.
func (c unbufferedPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, err := c.Read(b)
	return n, c.RemoteAddr(), err
}

// WriteTo implements net.PacketConn.WriteTo.
//
// addr is ignored.
func (c unbufferedPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return c.Write(b)
}

// NewUnbufferedPair creates a connected pair (socketpair) of unbuffered net.Conns.
//
// The implementation is inherently message oriented due to the lack of buffering. For each WriteXxx call, there is one opportunity to read the written data. If the buffer provided to ReadXxx is smaller than the buffer provided to the corresponding WriteXxx call, the unreadable data will be silently discarded.
func NewUnbufferedPair() (net.Conn, net.Conn) {
	c1, c2 := net.Pipe()
	return &unbufferedPacketConn{c1}, &unbufferedPacketConn{c2}
}

var errClosed = errors.New("closed")

// unbufferedAddr is similar to net.pipeAddr.
type unbufferedAddr struct{}

func (unbufferedAddr) Network() string { return "unbuffered" }

func (unbufferedAddr) String() string { return "unbuffered" }

// Unbuffered is an unbuffered net.Listener implemented in memory in userspace.
//
// Lighter weight, but not compatible with all applications. Creates unbuffered net.Conns.
type Unbuffered struct {
	acceptQueue chan net.Conn
	cancel      chan struct{}
}

var _ net.Listener = (*Unbuffered)(nil)

// NewUnbuffered creates a new unbuffered listener.
func NewUnbuffered() *Unbuffered {
	return &Unbuffered{make(chan net.Conn), make(chan struct{})}
}

// Dial creates a new unbuffered connection.
func (u *Unbuffered) Dial(ctx context.Context) (net.Conn, error) {
	// Check for closure/cancellation first.
	select {
	case <-u.cancel:
		return nil, errClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c1, c2 := NewUnbufferedPair()

	select {
	case <-u.cancel:
		c1.Close()
		c2.Close()
		return nil, errClosed
	case <-ctx.Done():
		c1.Close()
		c2.Close()
		return nil, ctx.Err()
	case u.acceptQueue <- c1:
		return c2, nil
	}
}

// Close implements net.Listener.Close.
func (u *Unbuffered) Close() error {
	close(u.cancel)
	return nil
}

// Accept implements net.Listener.Accept.
func (u *Unbuffered) Accept() (net.Conn, error) {
	// Check for closure first.
	select {
	case <-u.cancel:
		return nil, errClosed
	default:
	}

	select {
	case <-u.cancel:
		return nil, errClosed
	case a := <-u.acceptQueue:
		return a, nil
	}
}

// Addr implements net.Listener.Addr.
func (u *Unbuffered) Addr() net.Addr {
	return unbufferedAddr{}
}
