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

// HasUpstream reports whether the current branch of the repository rooted at
// path has an upstream tracking branch configured. It relies on the exit code
// of "git rev-parse @{u}" rather than parsing output, so it is not fooled by an
// empty capture from an unrelated failure.
func HasUpstream(path string) bool {
	params := []string{"-C", path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
	_, err := execute.ExternalProgramCapture("git", params...)

	return err == nil
}
