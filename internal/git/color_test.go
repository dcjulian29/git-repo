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

package git

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

const (
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
)

// forceColor enables ANSI output for the duration of a test regardless of
// whether the test output is a terminal, so color choices can be asserted.
func forceColor(t *testing.T) {
	t.Helper()

	previous := color.NoColor
	color.NoColor = false

	t.Cleanup(func() { color.NoColor = previous })
}

func TestColorBool(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name     string
		value    bool
		redOnYes bool
		text     string
		ansi     string
	}{
		{"false is green no", false, true, "no", ansiGreen},
		{"true red on yes", true, true, "yes", ansiRed},
		{"true yellow on yes", true, false, "yes", ansiYellow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColorBool(tt.value, tt.redOnYes)

			if !strings.Contains(got, tt.text) {
				t.Fatalf("expected text %q in %q", tt.text, got)
			}

			if !strings.Contains(got, tt.ansi) {
				t.Fatalf("expected color %q in %q", tt.ansi, got)
			}
		})
	}
}

func TestActionLabel(t *testing.T) {
	forceColor(t)

	if got := ActionLabel(false, "dirty", true); got != "" {
		t.Fatalf("expected empty label when flag is false, got %q", got)
	}

	if got := ActionLabel(true, "diverged", true); !strings.Contains(got, ansiRed) {
		t.Fatalf("expected red label, got %q", got)
	}

	if got := ActionLabel(true, "no upstream", false); !strings.Contains(got, ansiYellow) {
		t.Fatalf("expected yellow label, got %q", got)
	}
}

func TestColorPath(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name   string
		status RepoStatus
		ansi   string
	}{
		{"clean is green", RepoStatus{Folder: "/r"}, ansiGreen},
		{"dirty is yellow", RepoStatus{Folder: "/r", Dirty: true}, ansiYellow},
		{"no upstream is yellow", RepoStatus{Folder: "/r", NoUpstream: true}, ansiYellow},
		{"push needed is red", RepoStatus{Folder: "/r", PushNeeded: true}, ansiRed},
		{"pull needed is red", RepoStatus{Folder: "/r", PullNeeded: true}, ansiRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ColorPath(tt.status); !strings.Contains(got, tt.ansi) {
				t.Fatalf("expected color %q in %q", tt.ansi, got)
			}
		})
	}
}
