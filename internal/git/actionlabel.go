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

import "github.com/fatih/color"

// ActionLabel converts a boolean status flag into a human-readable coloured
// label (e.g. "dirty", "push needed") or an empty string when f is false.
// When redOnYes is true the label is rendered in red; otherwise yellow.
func ActionLabel(f bool, label string, redOnYes bool) string {
	if !f {
		return ""
	}

	if redOnYes {
		return color.RedString(label)
	}

	return color.YellowString(label)
}
