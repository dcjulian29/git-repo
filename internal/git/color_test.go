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
		text     string
		ansi     string
		value    bool
		redOnYes bool
	}{
		{name: "false is green no", value: false, redOnYes: true, text: "no", ansi: ansiGreen},
		{name: "true red on yes", value: true, redOnYes: true, text: "yes", ansi: ansiRed},
		{name: "true yellow on yes", value: true, redOnYes: false, text: "yes", ansi: ansiYellow},
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

func TestColorPath(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name   string
		ansi   string
		status RepoStatus
	}{
		{name: "clean is green", status: RepoStatus{Folder: "/r"}, ansi: ansiGreen},
		{name: "dirty is yellow", status: RepoStatus{Folder: "/r", Dirty: true}, ansi: ansiYellow},
		{name: "no upstream is yellow", status: RepoStatus{Folder: "/r", NoUpstream: true}, ansi: ansiYellow},
		{name: "push needed is red", status: RepoStatus{Folder: "/r", PushNeeded: true}, ansi: ansiRed},
		{name: "pull needed is red", status: RepoStatus{Folder: "/r", PullNeeded: true}, ansi: ansiRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ColorPath(tt.status); !strings.Contains(got, tt.ansi) {
				t.Fatalf("expected color %q in %q", tt.ansi, got)
			}
		})
	}
}
