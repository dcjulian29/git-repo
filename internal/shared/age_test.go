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
		want string
		d    time.Duration
	}{
		{name: "negative clamps to zero", d: -5 * time.Hour, want: "0m"},
		{name: "minutes", d: 45 * time.Minute, want: "45m"},
		{name: "hours", d: 5 * time.Hour, want: "5h"},
		{name: "days", d: 3 * ageDay, want: "3d"},
		{name: "months only", d: 5 * ageMonth, want: "5mo"},
		{name: "months and days", d: 5*ageMonth + 12*ageDay, want: "5mo12d"},
		{name: "years only", d: 2 * ageYear, want: "2y"},
		{name: "years and months", d: ageYear + 2*ageMonth, want: "1y2mo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Age(tt.d); got != tt.want {
				t.Fatalf("Age(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
