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

// Package memnet provides portable in-memory userspace "net" package implementations.
//
// Allows portable and hermetic testing and connecting. Useful for all cases where a Unix socket would be used, but you don't actually need to cross process boundaries.
//
// Buffered implementations are based on gVisor's Unix socket implementation.
package memnet
