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

package status

import (
	"testing"

	"github.com/dcjulian29/git-repo/internal/git"
)

// resetFilters clears every filter flag so each case starts from a known state.
func resetFilters() {
	actions = false
	dirty = false
	push = false
	pull = false
	diverged = false
	untracked = false
	noUpstream = false
}

func TestFilterDefaultShowsEverything(t *testing.T) {
	resetFilters()

	if !filter(git.RepoStatus{}) {
		t.Fatal("with no flags set, a clean repository should be shown")
	}
}

func TestFilterIndividualFlags(t *testing.T) {
	tests := []struct {
		name   string
		set    func()
		status git.RepoStatus
	}{
		{"dirty", func() { dirty = true }, git.RepoStatus{Dirty: true}},
		{"push", func() { push = true }, git.RepoStatus{PushNeeded: true}},
		{"pull", func() { pull = true }, git.RepoStatus{PullNeeded: true}},
		{"diverged", func() { diverged = true }, git.RepoStatus{Diverged: true}},
		{"untracked", func() { untracked = true }, git.RepoStatus{Untracked: true}},
		{"no upstream", func() { noUpstream = true }, git.RepoStatus{NoUpstream: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFilters()
			tt.set()

			if !filter(tt.status) {
				t.Fatalf("%s filter should match a matching repository", tt.name)
			}

			if filter(git.RepoStatus{}) {
				t.Fatalf("%s filter should exclude a clean repository", tt.name)
			}
		})
	}
}

func TestFilterActionsIncludesNoUpstream(t *testing.T) {
	resetFilters()
	actions = true

	if !filter(git.RepoStatus{NoUpstream: true}) {
		t.Fatal("--actions should include a no-upstream repository")
	}

	if filter(git.RepoStatus{}) {
		t.Fatal("--actions should exclude a clean repository")
	}
}
