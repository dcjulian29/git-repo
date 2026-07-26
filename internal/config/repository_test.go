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

package config

import "testing"

func TestRepositoryManaged(t *testing.T) {
	enabled := true
	disabled := false

	tests := []struct {
		repo Repository
		name string
		want bool
	}{
		{name: "unset defaults to managed", repo: Repository{Name: "a"}, want: true},
		{name: "explicit true is managed", repo: Repository{Name: "a", Manage: &enabled}, want: true},
		{name: "explicit false is skipped", repo: Repository{Name: "a", Manage: &disabled}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.repo.Managed(); got != tt.want {
				t.Fatalf("Managed() = %v, want %v", got, tt.want)
			}
		})
	}
}
