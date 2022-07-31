// Copyright 2022 The gVisor Authors.
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
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestUnbuffered_netHTTP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const want = "hello world"

	l := NewUnbuffered()

	server := http.Server{Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte(want)); err != nil {
			t.Error("http.ResponseWriter.Write:", err)
		}
	})}

	var eg errgroup.Group
	eg.Go(func() error {
		return server.Serve(l)
	})

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return l.DialContext(ctx)
	}
	client := http.Client{Transport: transport}

	response, err := client.Get(fmt.Sprintf("http://%s/", l.Addr()))
	if err != nil {
		t.Fatal("http.Client.Get:", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("got http.Response.StatusCode = %d, want = %d", response.StatusCode, http.StatusOK)
	}
	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("io.ReadAll(http.Response.Body):", err)
	}
	if got := string(got); got != want {
		t.Errorf("got http.Response.Body = %q, want = %q", got, want)
	}

	if err := server.Shutdown(ctx); err != nil {
		t.Fatal("http.Server.Shutdown:", err)
	}

	if err := eg.Wait(); err != http.ErrServerClosed {
		t.Errorf("got http.Server.Serve = %v, want = %v", err, http.ErrServerClosed)
	}
}
