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

import "github.com/dcjulian29/go-toolbox/execute"

// CaptureOutput runs git with args inside dir and returns trimmed stdout.
// On non-zero exit an empty string is returned (errors are intentionally
// swallowed so callers can use the output as a boolean-like signal).
func CaptureOutput(path string, args ...string) string {
	params := append([]string{"-C", path}, args...)
	out, _ := execute.ExternalProgramCapture("git", params...)

	return out
}
