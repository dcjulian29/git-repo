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

type (
	// Target pairs a caller-facing display name with the GitHub repository to query.
	Target struct {
		// Name is the display name for the repository (its configuration name).
		Name string

		// Repo is the parsed owner and name used to build API requests.
		Repo Repo
	}

	// Result holds the open items fetched for a single target, or the error
	// encountered while fetching them.
	Result struct {
		Err    error
		Target Target
		Items  []Item
	}
)

// Gather fetches items matching opts for every target concurrently, bounded by
// limit, and returns one Result per target in input order.
func (c *Client) Gather(ctx context.Context, targets []Target, limit int, opts ListOptions) []Result {
	return shared.ParallelMap(targets, limit, func(t Target) Result {
		items, err := c.ListItems(ctx, t.Repo, opts)

		return Result{Target: t, Items: items, Err: err}
	})
}
