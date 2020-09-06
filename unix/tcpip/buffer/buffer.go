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

// Package buffer provides the implementation of a buffer view.
package buffer

import (
	"bytes"
)

// View is a slice of a buffer, with convenience methods.
type View []byte

// NewView allocates a new buffer and returns an initialized view that covers
// the whole buffer.
func NewView(size int) View {
	return make(View, size)
}

// NewViewFromBytes allocates a new buffer and copies in the given bytes.
func NewViewFromBytes(b []byte) View {
	return append(View(nil), b...)
}

// TrimFront removes the first "count" bytes from the visible section of the
// buffer.
func (v *View) TrimFront(count int) {
	*v = (*v)[count:]
}

// CapLength irreversibly reduces the length of the visible section of the
// buffer to the value specified.
func (v *View) CapLength(length int) {
	// We also set the slice cap because if we don't, one would be able to
	// expand the view back to include the region just excluded. We want to
	// prevent that to avoid potential data leak if we have uninitialized
	// data in excluded region.
	*v = (*v)[:length:length]
}

// Reader returns a bytes.Reader for v.
func (v *View) Reader() bytes.Reader {
	var r bytes.Reader
	r.Reset(*v)
	return r
}

// IsEmpty returns whether v is of length zero.
func (v View) IsEmpty() bool {
	return len(v) == 0
}

// Size returns the length of v.
func (v View) Size() int {
	return len(v)
}
