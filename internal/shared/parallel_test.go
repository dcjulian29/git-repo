/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shared

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestParallelMapOrder verifies that results are returned in input order
// regardless of the order in which the workers complete.
func TestParallelMapOrder(t *testing.T) {
	inputs := make([]string, 20)
	for i := range inputs {
		inputs[i] = strconv.Itoa(i)
	}

	out := ParallelMap(inputs, 4, func(s string) string {
		return s + "!"
	})

	if len(out) != len(inputs) {
		t.Fatalf("got %d results, want %d", len(out), len(inputs))
	}

	for i, v := range out {
		if want := inputs[i] + "!"; v != want {
			t.Fatalf("order broken at index %d: got %q, want %q", i, v, want)
		}
	}
}

// TestParallelMapRespectsLimit verifies that no more than limit invocations of
// fn run concurrently.
func TestParallelMapRespectsLimit(t *testing.T) {
	const limit = 4

	inputs := make([]string, 20)
	for i := range inputs {
		inputs[i] = strconv.Itoa(i)
	}

	var (
		inflight    atomic.Int64
		maxInflight int64
		mu          sync.Mutex
	)

	ParallelMap(inputs, limit, func(s string) string {
		n := inflight.Add(1)

		mu.Lock()
		if n > maxInflight {
			maxInflight = n
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)
		inflight.Add(-1)

		return s
	})

	if maxInflight > limit {
		t.Fatalf("max concurrent invocations %d exceeded limit %d", maxInflight, limit)
	}
}

// TestParallelMapLimitFloor verifies that a limit below one is treated as one
// rather than deadlocking or panicking.
func TestParallelMapLimitFloor(t *testing.T) {
	out := ParallelMap([]string{"a", "b"}, 0, func(s string) string {
		return s
	})

	if len(out) != 2 || out[0] != "a" || out[1] != "b" {
		t.Fatalf("unexpected result with zero limit: %v", out)
	}
}

// TestParallelMapEmpty verifies that an empty input slice returns an empty
// result without spawning work.
func TestParallelMapEmpty(t *testing.T) {
	out := ParallelMap(nil, 4, func(s string) string {
		return s
	})

	if len(out) != 0 {
		t.Fatalf("expected no results, got %d", len(out))
	}
}

// TestDefaultConcurrencyBounds verifies the default stays within the intended
// range on any host.
func TestDefaultConcurrencyBounds(t *testing.T) {
	n := DefaultConcurrency()
	if n < 4 || n > 16 {
		t.Fatalf("DefaultConcurrency() = %d, want between 4 and 16", n)
	}
}
