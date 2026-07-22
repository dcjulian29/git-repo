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

import "testing"

func TestBrowserCommand(t *testing.T) {
	const url = "https://github.com/dcjulian29/git-repo/pull/1"

	tests := []struct {
		goos string
		name string
		last string
	}{
		{"windows", "rundll32", url},
		{"darwin", "open", url},
		{"linux", "xdg-open", url},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			name, args := browserCommand(tt.goos, url)

			if name != tt.name {
				t.Fatalf("launcher = %q, want %q", name, tt.name)
			}

			if len(args) == 0 || args[len(args)-1] != tt.last {
				t.Fatalf("args %v should end with the url %q", args, tt.last)
			}
		})
	}
}
