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

package goversion

import "testing"

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "trims prefix and newline", raw: "go1.26.5\n", want: "1.26.5"},
		{name: "trims surrounding spaces", raw: "  go1.25.0  ", want: "1.25.0"},
		{name: "already normalized", raw: "1.26.5", want: "1.26.5"},
		{name: "empty stays empty", raw: "\n", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeVersion(tt.raw); got != tt.want {
				t.Fatalf("normalizeVersion(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
