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

// Package syserr contains sandbox-internal errors. These errors are distinct
// from both the errors returned by host system calls and the errors returned
// to sandboxed applications.
package syserr

// Error represents an internal error.
type Error struct {
	// message is the human readable form of this Error.
	message string
}

// NewWithoutTranslation creates a new Error. If translation is attempted on
// the error, translation will fail.
//
// NewWithoutTranslation may be called at any time, but static errors should
// be declared as global variables and dynamic errors should be used sparingly.
func NewWithoutTranslation(message string) *Error {
	return &Error{message}
}

// String implements fmt.Stringer.String.
func (e *Error) String() string {
	if e == nil {
		return "<nil>"
	}
	return e.message
}

var (
	ErrNotPermitted               = NewWithoutTranslation("operation not permitted")
	ErrNoFileOrDir                = NewWithoutTranslation("no such file or directory")
	ErrNoProcess                  = NewWithoutTranslation("no such process")
	ErrInterrupted                = NewWithoutTranslation("interrupted system call")
	ErrIO                         = NewWithoutTranslation("I/O error")
	ErrDeviceOrAddress            = NewWithoutTranslation("no such device or address")
	ErrTooManyArgs                = NewWithoutTranslation("argument list too long")
	ErrEcec                       = NewWithoutTranslation("exec format error")
	ErrBadFD                      = NewWithoutTranslation("bad file number")
	ErrNoChild                    = NewWithoutTranslation("no child processes")
	ErrTryAgain                   = NewWithoutTranslation("try again")
	ErrNoMemory                   = NewWithoutTranslation("out of memory")
	ErrPermissionDenied           = NewWithoutTranslation("permission denied")
	ErrBadAddress                 = NewWithoutTranslation("bad address")
	ErrNotBlockDevice             = NewWithoutTranslation("block device required")
	ErrBusy                       = NewWithoutTranslation("device or resource busy")
	ErrExists                     = NewWithoutTranslation("file exists")
	ErrCrossDeviceLink            = NewWithoutTranslation("cross-device link")
	ErrNoDevice                   = NewWithoutTranslation("no such device")
	ErrNotDir                     = NewWithoutTranslation("not a directory")
	ErrIsDir                      = NewWithoutTranslation("is a directory")
	ErrInvalidArgument            = NewWithoutTranslation("invalid argument")
	ErrFileTableOverflow          = NewWithoutTranslation("file table overflow")
	ErrTooManyOpenFiles           = NewWithoutTranslation("too many open files")
	ErrNotTTY                     = NewWithoutTranslation("not a typewriter")
	ErrTestFileBusy               = NewWithoutTranslation("text file busy")
	ErrFileTooBig                 = NewWithoutTranslation("file too large")
	ErrNoSpace                    = NewWithoutTranslation("no space left on device")
	ErrIllegalSeek                = NewWithoutTranslation("illegal seek")
	ErrReadOnlyFS                 = NewWithoutTranslation("read-only file system")
	ErrTooManyLinks               = NewWithoutTranslation("too many links")
	ErrBrokenPipe                 = NewWithoutTranslation("broken pipe")
	ErrDomain                     = NewWithoutTranslation("math argument out of domain of func")
	ErrRange                      = NewWithoutTranslation("math result not representable")
	ErrDeadlock                   = NewWithoutTranslation("resource deadlock would occur")
	ErrNameTooLong                = NewWithoutTranslation("file name too long")
	ErrNoLocksAvailable           = NewWithoutTranslation("no record locks available")
	ErrInvalidSyscall             = NewWithoutTranslation("invalid system call number")
	ErrDirNotEmpty                = NewWithoutTranslation("directory not empty")
	ErrLinkLoop                   = NewWithoutTranslation("too many symbolic links encountered")
	ErrNoMessage                  = NewWithoutTranslation("no message of desired type")
	ErrIdentifierRemoved          = NewWithoutTranslation("identifier removed")
	ErrChannelOutOfRange          = NewWithoutTranslation("channel number out of range")
	ErrLevelTwoNotSynced          = NewWithoutTranslation("level 2 not synchronized")
	ErrLevelThreeHalted           = NewWithoutTranslation("level 3 halted")
	ErrLevelThreeReset            = NewWithoutTranslation("level 3 reset")
	ErrLinkNumberOutOfRange       = NewWithoutTranslation("link number out of range")
	ErrProtocolDriverNotAttached  = NewWithoutTranslation("protocol driver not attached")
	ErrNoCSIAvailable             = NewWithoutTranslation("no CSI structure available")
	ErrLevelTwoHalted             = NewWithoutTranslation("level 2 halted")
	ErrInvalidExchange            = NewWithoutTranslation("invalid exchange")
	ErrInvalidRequestDescriptor   = NewWithoutTranslation("invalid request descriptor")
	ErrExchangeFull               = NewWithoutTranslation("exchange full")
	ErrNoAnode                    = NewWithoutTranslation("no anode")
	ErrInvalidRequestCode         = NewWithoutTranslation("invalid request code")
	ErrInvalidSlot                = NewWithoutTranslation("invalid slot")
	ErrBadFontFile                = NewWithoutTranslation("bad font file format")
	ErrNotStream                  = NewWithoutTranslation("device not a stream")
	ErrNoDataAvailable            = NewWithoutTranslation("no data available")
	ErrTimerExpired               = NewWithoutTranslation("timer expired")
	ErrStreamsResourceDepleted    = NewWithoutTranslation("out of streams resources")
	ErrMachineNotOnNetwork        = NewWithoutTranslation("machine is not on the network")
	ErrPackageNotInstalled        = NewWithoutTranslation("package not installed")
	ErrIsRemote                   = NewWithoutTranslation("object is remote")
	ErrNoLink                     = NewWithoutTranslation("link has been severed")
	ErrAdvertise                  = NewWithoutTranslation("advertise error")
	ErrSRMount                    = NewWithoutTranslation("srmount error")
	ErrSendCommunication          = NewWithoutTranslation("communication error on send")
	ErrProtocol                   = NewWithoutTranslation("protocol error")
	ErrMultihopAttempted          = NewWithoutTranslation("multihop attempted")
	ErrRFS                        = NewWithoutTranslation("RFS specific error")
	ErrInvalidDataMessage         = NewWithoutTranslation("not a data message")
	ErrOverflow                   = NewWithoutTranslation("value too large for defined data type")
	ErrNetworkNameNotUnique       = NewWithoutTranslation("name not unique on network")
	ErrFDInBadState               = NewWithoutTranslation("file descriptor in bad state")
	ErrRemoteAddressChanged       = NewWithoutTranslation("remote address changed")
	ErrSharedLibraryInaccessible  = NewWithoutTranslation("can not access a needed shared library")
	ErrCorruptedSharedLibrary     = NewWithoutTranslation("accessing a corrupted shared library")
	ErrLibSectionCorrupted        = NewWithoutTranslation(".lib section in a.out corrupted")
	ErrTooManySharedLibraries     = NewWithoutTranslation("attempting to link in too many shared libraries")
	ErrSharedLibraryExeced        = NewWithoutTranslation("cannot exec a shared library directly")
	ErrIllegalByteSequence        = NewWithoutTranslation("illegal byte sequence")
	ErrShouldRestart              = NewWithoutTranslation("interrupted system call should be restarted")
	ErrStreamPipe                 = NewWithoutTranslation("streams pipe error")
	ErrTooManyUsers               = NewWithoutTranslation("too many users")
	ErrNotASocket                 = NewWithoutTranslation("socket operation on non-socket")
	ErrDestinationAddressRequired = NewWithoutTranslation("destination address required")
	ErrMessageTooLong             = NewWithoutTranslation("message too long")
	ErrWrongProtocolForSocket     = NewWithoutTranslation("protocol wrong type for socket")
	ErrProtocolNotAvailable       = NewWithoutTranslation("protocol not available")
	ErrProtocolNotSupported       = NewWithoutTranslation("protocol not supported")
	ErrSocketNotSupported         = NewWithoutTranslation("socket type not supported")
	ErrEndpointOperation          = NewWithoutTranslation("operation not supported on transport endpoint")
	ErrProtocolFamilyNotSupported = NewWithoutTranslation("protocol family not supported")
	ErrAddressFamilyNotSupported  = NewWithoutTranslation("address family not supported by protocol")
	ErrAddressInUse               = NewWithoutTranslation("address already in use")
	ErrAddressNotAvailable        = NewWithoutTranslation("cannot assign requested address")
	ErrNetworkDown                = NewWithoutTranslation("network is down")
	ErrNetworkUnreachable         = NewWithoutTranslation("network is unreachable")
	ErrNetworkReset               = NewWithoutTranslation("network dropped connection because of reset")
	ErrConnectionAborted          = NewWithoutTranslation("software caused connection abort")
	ErrConnectionReset            = NewWithoutTranslation("connection reset by peer")
	ErrNoBufferSpace              = NewWithoutTranslation("no buffer space available")
	ErrAlreadyConnected           = NewWithoutTranslation("transport endpoint is already connected")
	ErrNotConnected               = NewWithoutTranslation("transport endpoint is not connected")
	ErrShutdown                   = NewWithoutTranslation("cannot send after transport endpoint shutdown")
	ErrTooManyRefs                = NewWithoutTranslation("too many references: cannot splice")
	ErrTimedOut                   = NewWithoutTranslation("connection timed out")
	ErrConnectionRefused          = NewWithoutTranslation("connection refused")
	ErrHostDown                   = NewWithoutTranslation("host is down")
	ErrNoRoute                    = NewWithoutTranslation("no route to host")
	ErrAlreadyInProgress          = NewWithoutTranslation("operation already in progress")
	ErrInProgress                 = NewWithoutTranslation("operation now in progress")
	ErrStaleFileHandle            = NewWithoutTranslation("stale file handle")
	ErrStructureNeedsCleaning     = NewWithoutTranslation("structure needs cleaning")
	ErrIsNamedFile                = NewWithoutTranslation("is a named type file")
	ErrRemoteIO                   = NewWithoutTranslation("remote I/O error")
	ErrQuotaExceeded              = NewWithoutTranslation("quota exceeded")
	ErrNoMedium                   = NewWithoutTranslation("no medium found")
	ErrWrongMediumType            = NewWithoutTranslation("wrong medium type")
	ErrCanceled                   = NewWithoutTranslation("operation canceled")
	ErrNoKey                      = NewWithoutTranslation("required key not available")
	ErrKeyExpired                 = NewWithoutTranslation("key has expired")
	ErrKeyRevoked                 = NewWithoutTranslation("key has been revoked")
	ErrKeyRejected                = NewWithoutTranslation("key was rejected by service")
	ErrOwnerDied                  = NewWithoutTranslation("owner died")
	ErrNotRecoverable             = NewWithoutTranslation("state not recoverable")

	// ErrWouldBlock translates to EWOULDBLOCK which is the same as EAGAIN
	// on Linux.
	ErrWouldBlock = NewWithoutTranslation("operation would block")
)
