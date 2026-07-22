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

package github

import (
	"context"

	"github.com/dcjulian29/git-repo/internal/shared"
)

// Target pairs a caller-facing display name with the GitHub repository to query.
type Target struct {
	// Name is the display name for the repository (its configuration name).
	Name string

	// Repo is the parsed owner and name used to build API requests.
	Repo Repo
}

// Result holds the open items fetched for a single target, or the error
// encountered while fetching them.
type Result struct {
	// Target is the repository the result is for.
	Target Target

	// Items are the open issues and pull requests, when Err is nil.
	Items []Item

	// Err is non-nil when the repository could not be fetched.
	Err error
}

// Gather fetches open items for every target concurrently, bounded by limit,
// and returns one Result per target in input order.
func (c *Client) Gather(ctx context.Context, targets []Target, limit int) []Result {
	return shared.ParallelMap(targets, limit, func(t Target) Result {
		items, err := c.OpenItems(ctx, t.Repo)

		return Result{Target: t, Items: items, Err: err}
	})
}
