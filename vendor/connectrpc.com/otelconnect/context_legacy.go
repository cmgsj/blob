// Copyright 2022-2024 The Connect Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !go1.21

package otelconnect

import (
	"context"
	"sync/atomic"
)

// afterFunc is a simple imitation of context.AfterFunc from Go 1.21.
// It is not as efficient as the real implementation, but it is sufficient
// for our purposes.
func afterFunc(ctx context.Context, f func()) (stop func() bool) {
	ctx, cancel := context.WithCancel(ctx)
	var once atomic.Bool
	go func() {
		<-ctx.Done()
		if once.CompareAndSwap(false, true) {
			f()
		}
	}()
	return func() bool {
		didStop := once.CompareAndSwap(false, true)
		if didStop {
			cancel()
		}
		return didStop
	}
}
