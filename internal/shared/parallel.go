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
	"runtime"
	"sync"
)

// DefaultConcurrency returns a sensible upper bound on the number of
// repositories processed at once. The per-repository git operations are
// network/IO bound, so more workers than CPUs improves throughput, but the
// total is capped to avoid overwhelming the host or tripping remote rate
// limits when many repositories are managed.
func DefaultConcurrency() int {
	n := runtime.NumCPU() * 2

	if n < 4 {
		n = 4
	}

	if n > 16 {
		n = 16
	}

	return n
}

// ParallelMap applies fn to every input concurrently, running at most limit
// invocations at a time, and returns the results in input order. A limit below
// one is treated as one. It is used to fan out per-repository work without
// spawning an unbounded number of git subprocesses.
func ParallelMap[T any](inputs []string, limit int, fn func(string) T) []T {
	if limit < 1 {
		limit = 1
	}

	results := make([]T, len(inputs))

	var wg sync.WaitGroup

	sem := make(chan struct{}, limit)

	for i, in := range inputs {
		wg.Add(1)

		go func(i int, in string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			results[i] = fn(in)
		}(i, in)
	}

	wg.Wait()

	return results
}
