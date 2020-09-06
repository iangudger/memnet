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

// Package tcpip provides the interfaces and related types that users of the
// tcpip stack will use in order to create endpoints used to send and receive
// data over the network stack.
package tcpip

// Error represents an error in the netstack error space. Using a special type
// ensures that errors outside of this space are not accidentally introduced.
//
// Note: to support save / restore, it is important that all tcpip errors have
// distinct error messages.
type Error struct {
	msg string

	ignoreStats bool
}

// String implements fmt.Stringer.String.
func (e *Error) String() string {
	if e == nil {
		return "<nil>"
	}
	return e.msg
}

// IgnoreStats indicates whether this error type should be included in failure
// counts in tcpip.Stats structs.
func (e *Error) IgnoreStats() bool {
	return e.ignoreStats
}

// Errors that can be returned by the network stack.
var (
	ErrUnknownProtocol           = &Error{msg: "unknown protocol"}
	ErrUnknownNICID              = &Error{msg: "unknown nic id"}
	ErrUnknownDevice             = &Error{msg: "unknown device"}
	ErrUnknownProtocolOption     = &Error{msg: "unknown option for protocol"}
	ErrDuplicateNICID            = &Error{msg: "duplicate nic id"}
	ErrDuplicateAddress          = &Error{msg: "duplicate address"}
	ErrNoRoute                   = &Error{msg: "no route"}
	ErrBadLinkEndpoint           = &Error{msg: "bad link layer endpoint"}
	ErrAlreadyBound              = &Error{msg: "endpoint already bound", ignoreStats: true}
	ErrInvalidEndpointState      = &Error{msg: "endpoint is in invalid state"}
	ErrAlreadyConnecting         = &Error{msg: "endpoint is already connecting", ignoreStats: true}
	ErrAlreadyConnected          = &Error{msg: "endpoint is already connected", ignoreStats: true}
	ErrNoPortAvailable           = &Error{msg: "no ports are available"}
	ErrPortInUse                 = &Error{msg: "port is in use"}
	ErrBadLocalAddress           = &Error{msg: "bad local address"}
	ErrClosedForSend             = &Error{msg: "endpoint is closed for send"}
	ErrClosedForReceive          = &Error{msg: "endpoint is closed for receive"}
	ErrWouldBlock                = &Error{msg: "operation would block", ignoreStats: true}
	ErrConnectionRefused         = &Error{msg: "connection was refused"}
	ErrTimeout                   = &Error{msg: "operation timed out"}
	ErrAborted                   = &Error{msg: "operation aborted"}
	ErrConnectStarted            = &Error{msg: "connection attempt started", ignoreStats: true}
	ErrDestinationRequired       = &Error{msg: "destination address is required"}
	ErrNotSupported              = &Error{msg: "operation not supported"}
	ErrQueueSizeNotSupported     = &Error{msg: "queue size querying not supported"}
	ErrNotConnected              = &Error{msg: "endpoint not connected"}
	ErrConnectionReset           = &Error{msg: "connection reset by peer"}
	ErrConnectionAborted         = &Error{msg: "connection aborted"}
	ErrNoSuchFile                = &Error{msg: "no such file"}
	ErrInvalidOptionValue        = &Error{msg: "invalid option value specified"}
	ErrNoLinkAddress             = &Error{msg: "no remote link address"}
	ErrBadAddress                = &Error{msg: "bad address"}
	ErrNetworkUnreachable        = &Error{msg: "network is unreachable"}
	ErrMessageTooLong            = &Error{msg: "message too long"}
	ErrNoBufferSpace             = &Error{msg: "no buffer space available"}
	ErrBroadcastDisabled         = &Error{msg: "broadcast socket option disabled"}
	ErrNotPermitted              = &Error{msg: "operation not permitted"}
	ErrAddressFamilyNotSupported = &Error{msg: "address family not supported by protocol"}
)

// Address is a byte slice cast as a string that represents the address of a
// network node. Or, in the case of unix endpoints, it may represent a path.
type Address string

// NICID is a number that uniquely identifies a NIC.
type NICID int32

// ShutdownFlags represents flags that can be passed to the Shutdown() method
// of the Endpoint interface.
type ShutdownFlags int

// Values of the flags that can be passed to the Shutdown() method. They can
// be OR'ed together.
const (
	ShutdownRead ShutdownFlags = 1 << iota
	ShutdownWrite
)

// FullAddress represents a full transport node address, as required by the
// Connect() and Bind() methods.
//
// +stateify savable
type FullAddress struct {
	// NIC is the ID of the NIC this address refers to.
	//
	// This may not be used by all endpoint types.
	NIC NICID

	// Addr is the network or link layer address.
	Addr Address

	// Port is the transport port.
	//
	// This may not be used by all endpoint types.
	Port uint16
}

// SockOptBool represents socket options which values have the bool type.
type SockOptBool int

const (
	// BroadcastOption is used by SetSockOptBool/GetSockOptBool to specify
	// whether datagram sockets are allowed to send packets to a broadcast
	// address.
	BroadcastOption SockOptBool = iota

	// CorkOption is used by SetSockOptBool/GetSockOptBool to specify if
	// data should be held until segments are full by the TCP transport
	// protocol.
	CorkOption

	// DelayOption is used by SetSockOptBool/GetSockOptBool to specify if
	// data should be sent out immediately by the transport protocol. For
	// TCP, it determines if the Nagle algorithm is on or off.
	DelayOption

	// KeepaliveEnabledOption is used by SetSockOptBool/GetSockOptBool to
	// specify whether TCP keepalive is enabled for this socket.
	KeepaliveEnabledOption

	// MulticastLoopOption is used by SetSockOptBool/GetSockOptBool to
	// specify whether multicast packets sent over a non-loopback interface
	// will be looped back.
	MulticastLoopOption

	// NoChecksumOption is used by SetSockOptBool/GetSockOptBool to specify
	// whether UDP checksum is disabled for this socket.
	NoChecksumOption

	// PasscredOption is used by SetSockOptBool/GetSockOptBool to specify
	// whether SCM_CREDENTIALS socket control messages are enabled.
	//
	// Only supported on Unix sockets.
	PasscredOption

	// QuickAckOption is stubbed out in SetSockOptBool/GetSockOptBool.
	QuickAckOption

	// ReceiveTClassOption is used by SetSockOptBool/GetSockOptBool to
	// specify if the IPV6_TCLASS ancillary message is passed with incoming
	// packets.
	ReceiveTClassOption

	// ReceiveTOSOption is used by SetSockOptBool/GetSockOptBool to specify
	// if the TOS ancillary message is passed with incoming packets.
	ReceiveTOSOption

	// ReceiveIPPacketInfoOption is used by SetSockOptBool/GetSockOptBool to
	// specify if more inforamtion is provided with incoming packets such as
	// interface index and address.
	ReceiveIPPacketInfoOption

	// ReuseAddressOption is used by SetSockOptBool/GetSockOptBool to
	// specify whether Bind() should allow reuse of local address.
	ReuseAddressOption

	// ReusePortOption is used by SetSockOptBool/GetSockOptBool to permit
	// multiple sockets to be bound to an identical socket address.
	ReusePortOption

	// V6OnlyOption is used by SetSockOptBool/GetSockOptBool to specify
	// whether an IPv6 socket is to be restricted to sending and receiving
	// IPv6 packets only.
	V6OnlyOption

	// IPHdrIncludedOption is used by SetSockOpt to indicate for a raw
	// endpoint that all packets being written have an IP header and the
	// endpoint should not attach an IP header.
	IPHdrIncludedOption
)

// SockOptInt represents socket options which values have the int type.
type SockOptInt int

const (
	// KeepaliveCountOption is used by SetSockOptInt/GetSockOptInt to
	// specify the number of un-ACKed TCP keepalives that will be sent
	// before the connection is closed.
	KeepaliveCountOption SockOptInt = iota

	// IPv4TOSOption is used by SetSockOptInt/GetSockOptInt to specify TOS
	// for all subsequent outgoing IPv4 packets from the endpoint.
	IPv4TOSOption

	// IPv6TrafficClassOption is used by SetSockOptInt/GetSockOptInt to
	// specify TOS for all subsequent outgoing IPv6 packets from the
	// endpoint.
	IPv6TrafficClassOption

	// MaxSegOption is used by SetSockOptInt/GetSockOptInt to set/get the
	// current Maximum Segment Size(MSS) value as specified using the
	// TCP_MAXSEG option.
	MaxSegOption

	// MTUDiscoverOption is used to set/get the path MTU discovery setting.
	//
	// NOTE: Setting this option to any other value than PMTUDiscoveryDont
	// is not supported and will fail as such, and getting this option will
	// always return PMTUDiscoveryDont.
	MTUDiscoverOption

	// MulticastTTLOption is used by SetSockOptInt/GetSockOptInt to control
	// the default TTL value for multicast messages. The default is 1.
	MulticastTTLOption

	// ReceiveQueueSizeOption is used in GetSockOptInt to specify that the
	// number of unread bytes in the input buffer should be returned.
	ReceiveQueueSizeOption

	// SendBufferSizeOption is used by SetSockOptInt/GetSockOptInt to
	// specify the send buffer size option.
	SendBufferSizeOption

	// ReceiveBufferSizeOption is used by SetSockOptInt/GetSockOptInt to
	// specify the receive buffer size option.
	ReceiveBufferSizeOption

	// SendQueueSizeOption is used in GetSockOptInt to specify that the
	// number of unread bytes in the output buffer should be returned.
	SendQueueSizeOption

	// TTLOption is used by SetSockOptInt/GetSockOptInt to control the
	// default TTL/hop limit value for unicast messages. The default is
	// protocol specific.
	//
	// A zero value indicates the default.
	TTLOption

	// TCPSynCountOption is used by SetSockOptInt/GetSockOptInt to specify
	// the number of SYN retransmits that TCP should send before aborting
	// the attempt to connect. It cannot exceed 255.
	//
	// NOTE: This option is currently only stubbed out and is no-op.
	TCPSynCountOption

	// TCPWindowClampOption is used by SetSockOptInt/GetSockOptInt to bound
	// the size of the advertised window to this value.
	//
	// NOTE: This option is currently only stubed out and is a no-op
	TCPWindowClampOption
)

// ErrorOption is used in GetSockOpt to specify that the last error reported by
// the endpoint should be cleared and returned.
type ErrorOption struct{}
