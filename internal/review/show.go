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

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/fatih/color"
)

// newClient resolves a token and returns an authenticated GitHub client.
func newClient() (*github.Client, error) {
	token, err := github.Token()
	if err != nil {
		return nil, err
	}

	return github.NewClient(token), nil
}

// ShowPull fetches and prints the details of a single pull request: metadata,
// mergeability and conflicts, the dependabot compatibility score when present,
// and the state of its checks.
func ShowPull(ctx context.Context, ref Ref) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	_, err = describe(ctx, client, ref)

	return err
}

// ShowIssue fetches and prints the details of a single issue.
func ShowIssue(ctx context.Context, ref Ref) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	renderIssue(ref, issue)

	return nil
}

// describe fetches a pull request with its checks and compatibility score,
// renders the summary, and returns the pull request for further action.
func describe(ctx context.Context, client *github.Client, ref Ref) (github.PullRequest, error) {
	pull, err := client.GetPull(ctx, ref.Repo, ref.Number)
	if err != nil {
		return github.PullRequest{}, err
	}

	checks, err := client.Checks(ctx, ref.Repo, pull.HeadSHA)
	if err != nil {
		return github.PullRequest{}, err
	}

	compat := github.CompatibilityScore(ctx, pull.Body)

	renderPull(ref, pull, checks, compat)

	return pull, nil
}

func renderPull(ref Ref, pull github.PullRequest, checks []github.CheckRun, compat string) {
	fmt.Printf("%s %s\n",
		color.CyanString("%s#%d", ref.Name, pull.Number),
		pull.Title)

	fmt.Printf("  %-14s %s\n", "Author:", pull.Author)
	fmt.Printf("  %-14s %s\n", "State:", pullState(pull))
	fmt.Printf("  %-14s %s → %s\n", "Branch:", pull.HeadRef, pull.BaseRef)
	fmt.Printf("  %-14s %s\n", "Mergeable:", mergeable(pull))

	if compat != "" {
		fmt.Printf("  %-14s %s\n", "Compatibility:", compat)
	}

	renderChecks(checks)
}

func pullState(pull github.PullRequest) string {
	if pull.Draft {
		return pull.State + " " + color.YellowString("(draft)")
	}

	return pull.State
}

func mergeable(pull github.PullRequest) string {
	switch {
	case pull.HasConflicts():
		return color.RedString("no") + fmt.Sprintf(" (%s)", pull.MergeableState)
	case pull.Mergeable == nil:
		return color.YellowString("unknown") + " (still computing)"
	default:
		return color.GreenString("yes") + fmt.Sprintf(" (%s)", pull.MergeableState)
	}
}

func renderChecks(checks []github.CheckRun) {
	if len(checks) == 0 {
		fmt.Printf("  %-14s %s\n", "Checks:", "none")

		return
	}

	summaries := aggregateChecks(checks)

	fmt.Printf("  Checks (%d):\n", len(summaries))

	for _, summary := range summaries {
		name := summary.run.Name
		if summary.count > 1 {
			name = fmt.Sprintf("%s ×%d", name, summary.count)
		}

		fmt.Printf("    %-30s %s\n", name, checkStatus(summary.run))
	}
}

// checkSummary collapses the check runs that share a display name into a single
// worst-case run and the number of runs that used that name.
type checkSummary struct {
	run   github.CheckRun
	count int
}

// aggregateChecks groups check runs by name, keeping the worst-case run for each
// name and counting how many shared it, in first-seen order. GitHub returns a
// separate run for every leg of a matrix or repeated trigger, so the same name
// can appear several times; collapsing them keeps the summary readable.
func aggregateChecks(checks []github.CheckRun) []checkSummary {
	index := make(map[string]int, len(checks))
	summaries := make([]checkSummary, 0, len(checks))

	for _, check := range checks {
		if i, ok := index[check.Name]; ok {
			summaries[i].count++

			if checkSeverity(check) > checkSeverity(summaries[i].run) {
				summaries[i].run = check
			}

			continue
		}

		index[check.Name] = len(summaries)
		summaries = append(summaries, checkSummary{run: check, count: 1})
	}

	return summaries
}

// checkSeverity ranks a check run so the worst outcome wins when a name repeats:
// a failure outranks a still-running check, which outranks a success.
func checkSeverity(check github.CheckRun) int {
	switch {
	case check.Completed() && !check.Passed():
		return 2
	case !check.Completed():
		return 1
	default:
		return 0
	}
}

// checkStatus returns the check's result text colored green when it passed,
// red when it failed, and yellow while it is still running.
func checkStatus(check github.CheckRun) string {
	switch {
	case !check.Completed():
		return color.YellowString(check.Status)
	case check.Passed():
		return color.GreenString(check.Conclusion)
	default:
		return color.RedString(check.Conclusion)
	}
}
