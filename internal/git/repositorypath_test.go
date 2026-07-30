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
	"path/filepath"
	"testing"
)

func TestRepositoryPath(t *testing.T) {
	tests := []struct {
		name string
		repo string
		want string
	}{
		{name: "simple name", repo: "widget", want: filepath.Join("root", "widget")},
		{name: "nested name", repo: "ansible/roles/base", want: filepath.Join("root", "ansible", "roles", "base")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RepositoryPath("root", tt.repo); got != tt.want {
				t.Fatalf("RepositoryPath(root, %q) = %q, want %q", tt.repo, got, tt.want)
			}
		})
	}
}
