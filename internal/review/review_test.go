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
	"errors"
	"testing"
	"time"

	"github.com/dcjulian29/git-repo/internal/github"
)

func sampleReport() Report {
	base := time.Now()

	return Report{
		Results: []github.Result{
			{
				Target: github.Target{Name: "repoA"},
				Items: []github.Item{
					{Number: 1, IsPull: false, CreatedAt: base.Add(-2 * time.Hour)},
					{Number: 2, IsPull: true, CreatedAt: base.Add(-1 * time.Hour)},
				},
			},
			{Target: github.Target{Name: "repoB"}, Err: errors.New("boom")},
			{
				Target: github.Target{Name: "repoC"},
				Items: []github.Item{
					{Number: 3, IsPull: true, CreatedAt: base.Add(-3 * time.Hour)},
				},
			},
		},
	}
}

func TestReportPullRequestsSortedOldestFirst(t *testing.T) {
	prs := sampleReport().PullRequests()

	if len(prs) != 2 {
		t.Fatalf("got %d pull requests, want 2", len(prs))
	}

	// #3 (repoC, 3h old) should sort before #2 (repoA, 1h old).
	if prs[0].Number != 3 || prs[0].Repo != "repoC" {
		t.Fatalf("first PR = #%d in %s, want #3 in repoC", prs[0].Number, prs[0].Repo)
	}

	if prs[1].Number != 2 || prs[1].Repo != "repoA" {
		t.Fatalf("second PR = #%d in %s, want #2 in repoA", prs[1].Number, prs[1].Repo)
	}
}

func TestReportIssuesExcludePulls(t *testing.T) {
	issues := sampleReport().Issues()

	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1", len(issues))
	}

	if issues[0].Number != 1 || issues[0].IsPull {
		t.Fatalf("issue = %+v, want #1 non-pull", issues[0])
	}
}

func TestReportFailures(t *testing.T) {
	failures := sampleReport().Failures()

	if len(failures) != 1 || failures[0].Target.Name != "repoB" {
		t.Fatalf("failures = %+v, want one for repoB", failures)
	}
}
