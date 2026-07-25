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
	"testing"

	"github.com/dcjulian29/git-repo/internal/github"
)

func TestFilterPulls(t *testing.T) {
	pulls := []NamedItem{
		{Repo: "r", Item: github.Item{Number: 1, IsPull: true}},                          // open, non-draft
		{Repo: "r", Item: github.Item{Number: 2, IsPull: true, Draft: true}},             // open draft
		{Repo: "r", Item: github.Item{Number: 3, IsPull: true, Merged: true}},            // merged
		{Repo: "r", Item: github.Item{Number: 4, IsPull: true}},                          // closed, unmerged
	}

	tests := []struct {
		mode string
		want []int
	}{
		{"open", []int{1, 2, 3, 4}}, // no extra filtering; state query already scoped these
		{"draft", []int{2}},
		{"merged", []int{3}},
		{"closed", []int{1, 2, 4}}, // everything not merged
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			got := filterPulls(pulls, tt.mode)

			if len(got) != len(tt.want) {
				t.Fatalf("mode %q returned %d items, want %d", tt.mode, len(got), len(tt.want))
			}

			for i, number := range tt.want {
				if got[i].Number != number {
					t.Fatalf("mode %q item %d = #%d, want #%d", tt.mode, i, got[i].Number, number)
				}
			}
		})
	}
}

func TestIssuesEmptyMessage(t *testing.T) {
	cases := map[string]string{
		"":       "No open issues found.",
		"open":   "No open issues found.",
		"closed": "No closed issues found.",
		"all":    "No issues found.",
	}

	for state, want := range cases {
		if got := issuesEmptyMessage(state); got != want {
			t.Fatalf("issuesEmptyMessage(%q) = %q, want %q", state, got, want)
		}
	}
}
