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
	"testing"
	"time"
)

func TestAge(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"negative clamps to zero", -5 * time.Hour, "0m"},
		{"minutes", 45 * time.Minute, "45m"},
		{"hours", 5 * time.Hour, "5h"},
		{"days", 3 * ageDay, "3d"},
		{"months only", 5 * ageMonth, "5mo"},
		{"months and days", 5*ageMonth + 12*ageDay, "5mo12d"},
		{"years only", 2 * ageYear, "2y"},
		{"years and months", ageYear + 2*ageMonth, "1y2mo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Age(tt.d); got != tt.want {
				t.Fatalf("Age(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
