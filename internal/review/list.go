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

package review

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
)

// ListPulls prints the pull requests across managed repositories selected by
// mode, one of "open", "closed", "merged", or "draft". When mergeState is
// non-empty, only pull requests with that merge state (for example "clean") are
// shown.
func ListPulls(ctx context.Context, mode, mergeState string, asJSON bool) error {
	state := "open"
	if mode == "closed" || mode == "merged" {
		state = "closed"
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	report, err := Collect(ctx, client, github.ListOptions{State: state})
	if err != nil {
		return err
	}

	PrintWarnings(os.Stderr, report)

	pulls := filterPulls(report.PullRequests(), mode)

	// The listing endpoint does not report merge state; enrich the pull requests
	// where it is meaningful (open and draft) or when it is being filtered on.
	if mode == "open" || mode == "draft" || mergeState != "" {
		pulls = enrichMergeState(ctx, client, pulls)
	}

	if mergeState != "" {
		pulls = filterByMergeState(pulls, mergeState)
	}

	if asJSON {
		return RenderJSON(pulls)
	}

	return RenderTable(pulls, "PR", pullsEmptyMessage(mode, mergeState), true)
}

// filterByMergeState keeps only the pull requests whose merge state matches
// state (case-insensitively).
func filterByMergeState(pulls []NamedItem, state string) []NamedItem {
	filtered := make([]NamedItem, 0, len(pulls))

	for _, item := range pulls {
		if strings.EqualFold(item.MergeableState, state) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// pullsEmptyMessage returns the "nothing found" notice for the given mode and
// optional merge-state filter.
func pullsEmptyMessage(mode, mergeState string) string {
	if mergeState != "" {
		return fmt.Sprintf("No %s pull requests with merge state %q found.", mode, mergeState)
	}

	return fmt.Sprintf("No %s pull requests found.", mode)
}

// enrichMergeState fetches each pull request's merge state concurrently and
// returns the pull requests with MergeableState populated.
func enrichMergeState(ctx context.Context, client *github.Client, pulls []NamedItem) []NamedItem {
	return shared.ParallelMap(pulls, shared.DefaultConcurrency(), func(item NamedItem) NamedItem {
		if pull, err := client.GetPull(ctx, item.ghRepo, item.Number); err == nil {
			item.MergeableState = pull.MergeableState
		}

		return item
	})
}

// filterPulls narrows pull requests to those matching mode. The state query has
// already limited them to open or closed; this applies the draft/merged
// distinction that the API does not express directly.
func filterPulls(pulls []NamedItem, mode string) []NamedItem {
	var predicate func(NamedItem) bool

	switch mode {
	case "draft":
		predicate = func(item NamedItem) bool { return item.Draft }
	case "closed":
		predicate = func(item NamedItem) bool { return !item.Merged }
	case "merged":
		predicate = func(item NamedItem) bool { return item.Merged }
	default:
		return pulls
	}

	filtered := make([]NamedItem, 0, len(pulls))

	for _, item := range pulls {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// ListIssues prints the issues across managed repositories matching opts.
func ListIssues(ctx context.Context, opts github.ListOptions, asJSON bool) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	report, err := Collect(ctx, client, opts)
	if err != nil {
		return err
	}

	PrintWarnings(os.Stderr, report)

	issues := report.Issues()

	if asJSON {
		return RenderJSON(issues)
	}

	return RenderTable(issues, "ISSUE", issuesEmptyMessage(opts.State), false)
}

// issuesEmptyMessage returns the "nothing found" notice appropriate to the
// requested state.
func issuesEmptyMessage(state string) string {
	if state == "all" {
		return "No issues found."
	}

	if state == "" {
		state = "open"
	}

	return fmt.Sprintf("No %s issues found.", state)
}
