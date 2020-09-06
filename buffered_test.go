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
	"fmt"
	"net"
	"testing"

	"golang.org/x/net/nettest"
)

func TestNewBufferedPairConformance(t *testing.T) {
	for _, network := range []string{"unix", "unixgram", "unixpacket"} {
		t.Run(network, func(t *testing.T) {
			nettest.TestConn(t, func() (c1, c2 net.Conn, stop func(), err error) {
				c1, c2, err = NewBufferedPair(network)
				if err != nil {
					return nil, nil, nil, err
				}
				stop = func() {
					c1.Close()
					c2.Close()
				}
				return
			})
		})
	}
}

func TestBufferedListenConformance(t *testing.T) {
	for _, network := range []string{"unix", "unixpacket"} {
		t.Run(network, func(t *testing.T) {
			nettest.TestConn(t, func() (c1, c2 net.Conn, stop func(), err error) {
				l, err := BufferedListen(network, &net.UnixAddr{
					Name: "test",
					Net:  network,
				})
				if err != nil {
					return nil, nil, nil, fmt.Errorf("BufferedListen: %w", err)
				}

				c1, err = l.Dial(nil)
				if err != nil {
					l.Close()
					return nil, nil, nil, fmt.Errorf("DialTCP: %w", err)
				}

				c2, err = l.Accept()
				if err != nil {
					c1.Close()
					l.Close()
					return nil, nil, nil, fmt.Errorf("Accept: %w", err)
				}

				stop = func() {
					c1.Close()
					c2.Close()
					l.Close()
				}
				return
			})
		})
	}
}
