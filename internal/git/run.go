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

// Run executes git with args inside the repository at path, binding the child
// process to the host terminal so the user sees git's live output. It returns
// an error when git exits non-zero.
func Run(path string, args ...string) error {
	params := append([]string{"-C", path}, args...)

	return execute.ExternalProgram("git", params...)
}
