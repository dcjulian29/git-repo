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

func TestIif(t *testing.T) {
	if got := Iif(true, "a", "b"); got != "a" {
		t.Fatalf("Iif(true) = %q, want %q", got, "a")
	}

	if got := Iif(false, "a", "b"); got != "b" {
		t.Fatalf("Iif(false) = %q, want %q", got, "b")
	}
}
