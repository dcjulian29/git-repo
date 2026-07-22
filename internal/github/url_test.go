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

import "testing"

func TestParseRepoValid(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		owner string
		repo  string
	}{
		{"https with .git", "https://github.com/dcjulian29/git-repo.git", "dcjulian29", "git-repo"},
		{"https without .git", "https://github.com/dcjulian29/git-repo", "dcjulian29", "git-repo"},
		{"https trailing slash", "https://github.com/dcjulian29/git-repo/", "dcjulian29", "git-repo"},
		{"scp-style ssh", "git@github.com:dcjulian29/git-repo.git", "dcjulian29", "git-repo"},
		{"ssh url", "ssh://git@github.com/dcjulian29/git-repo.git", "dcjulian29", "git-repo"},
		{"surrounding spaces", "  https://github.com/dcjulian29/git-repo.git  ", "dcjulian29", "git-repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepo(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Owner != tt.owner || got.Name != tt.repo {
				t.Fatalf("ParseRepo(%q) = %+v, want owner=%q name=%q", tt.in, got, tt.owner, tt.repo)
			}

			if want := tt.owner + "/" + tt.repo; got.String() != want {
				t.Fatalf("String() = %q, want %q", got.String(), want)
			}
		})
	}
}

func TestParseRepoInvalid(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"gitlab host", "https://gitlab.com/foo/bar.git"},
		{"bitbucket ssh", "git@bitbucket.org:foo/bar.git"},
		{"missing name", "https://github.com/onlyowner"},
		{"extra path segment", "https://github.com/owner/repo/extra"},
		{"not a url", "just some text"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseRepo(tt.in); err == nil {
				t.Fatalf("ParseRepo(%q) expected an error, got nil", tt.in)
			}
		})
	}
}
