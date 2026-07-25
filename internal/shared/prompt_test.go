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

package shared

import (
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"trimmed answer", "  accepted  \n", "accepted"},
		{"empty line", "\n", ""},
		{"eof", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder

			got, err := Prompt(strings.NewReader(tt.input), &out, "label? ")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("Prompt(%q) = %q, want %q", tt.input, got, tt.want)
			}

			if !strings.Contains(out.String(), "label?") {
				t.Fatalf("question not written, got %q", out.String())
			}
		})
	}
}
