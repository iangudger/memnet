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

package syserr

import (
	"github.com/iangudger/memnet/unix/tcpip"
)

var (
	ErrUnknownProtocol       = NewWithoutTranslation(tcpip.ErrUnknownProtocol.String())
	ErrUnknownNICID          = NewWithoutTranslation(tcpip.ErrUnknownNICID.String())
	ErrUnknownDevice         = NewWithoutTranslation(tcpip.ErrUnknownDevice.String())
	ErrUnknownProtocolOption = NewWithoutTranslation(tcpip.ErrUnknownProtocolOption.String())
	ErrDuplicateNICID        = NewWithoutTranslation(tcpip.ErrDuplicateNICID.String())
	ErrDuplicateAddress      = NewWithoutTranslation(tcpip.ErrDuplicateAddress.String())
	ErrBadLinkEndpoint       = NewWithoutTranslation(tcpip.ErrBadLinkEndpoint.String())
	ErrAlreadyBound          = NewWithoutTranslation(tcpip.ErrAlreadyBound.String())
	ErrInvalidEndpointState  = NewWithoutTranslation(tcpip.ErrInvalidEndpointState.String())
	ErrAlreadyConnecting     = NewWithoutTranslation(tcpip.ErrAlreadyConnecting.String())
	ErrNoPortAvailable       = NewWithoutTranslation(tcpip.ErrNoPortAvailable.String())
	ErrPortInUse             = NewWithoutTranslation(tcpip.ErrPortInUse.String())
	ErrBadLocalAddress       = NewWithoutTranslation(tcpip.ErrBadLocalAddress.String())
	ErrClosedForSend         = NewWithoutTranslation(tcpip.ErrClosedForSend.String())
	ErrClosedForReceive      = NewWithoutTranslation(tcpip.ErrClosedForReceive.String())
	ErrTimeout               = NewWithoutTranslation(tcpip.ErrTimeout.String())
	ErrAborted               = NewWithoutTranslation(tcpip.ErrAborted.String())
	ErrConnectStarted        = NewWithoutTranslation(tcpip.ErrConnectStarted.String())
	ErrDestinationRequired   = NewWithoutTranslation(tcpip.ErrDestinationRequired.String())
	ErrNotSupported          = NewWithoutTranslation(tcpip.ErrNotSupported.String())
	ErrQueueSizeNotSupported = NewWithoutTranslation(tcpip.ErrQueueSizeNotSupported.String())
	ErrNoSuchFile            = NewWithoutTranslation(tcpip.ErrNoSuchFile.String())
	ErrInvalidOptionValue    = NewWithoutTranslation(tcpip.ErrInvalidOptionValue.String())
	ErrBroadcastDisabled     = NewWithoutTranslation(tcpip.ErrBroadcastDisabled.String())
	ErrNotPermittedNet       = NewWithoutTranslation(tcpip.ErrNotPermitted.String())
)
