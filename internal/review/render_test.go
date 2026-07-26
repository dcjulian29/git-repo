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
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/fatih/color"
)

func TestColorHandle(t *testing.T) {
	previous := color.NoColor
	color.NoColor = false

	t.Cleanup(func() { color.NoColor = previous })

	const (
		green  = "\x1b[32m"
		yellow = "\x1b[33m"
		red    = "\x1b[31m"
		cyan   = "\x1b[36m"
	)

	cases := map[string]string{
		"clean":    green,
		"unstable": red,
		"dirty":    red,
		"unknown":  yellow,
		"":         cyan,
	}

	for state, ansi := range cases {
		if got := colorHandle("git-repo#1", state); !strings.Contains(got, ansi) {
			t.Fatalf("colorHandle(state=%q) = %q, want colour %q", state, got, ansi)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{"fits", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"truncated with ellipsis", "hello world", 8, "hello w…"},
		{"single column", "hello", 1, "…"},
		{"zero column", "hello", 0, ""},
		{"unicode runes", "日本語テスト", 3, "日本…"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.in, tt.max)
			if got != tt.want {
				t.Fatalf("truncate(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
			}

			if tt.max > 0 && utf8.RuneCountInString(got) > tt.max {
				t.Fatalf("result %q exceeds max width %d", got, tt.max)
			}
		})
	}
}

func TestTitleBudget(t *testing.T) {
	// 120 - 13 chrome (4 cols) - 40 other = 67.
	if got := titleBudget(120, 40, 4); got != 67 {
		t.Fatalf("titleBudget wide = %d, want 67", got)
	}

	// Five columns add 3 more chrome: 120 - 16 - 40 = 64.
	if got := titleBudget(120, 40, 5); got != 64 {
		t.Fatalf("titleBudget five columns = %d, want 64", got)
	}

	// A narrow terminal clamps to the floor rather than going negative.
	if got := titleBudget(40, 40, 4); got != minTitleWidth {
		t.Fatalf("titleBudget narrow = %d, want %d", got, minTitleWidth)
	}
}
