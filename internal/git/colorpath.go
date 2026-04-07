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

// ColorPath returns an ANSI-coloured string for the repository path:
//   - yellow  when the working tree is dirty
//   - red     when a push or pull is needed
//   - green   when everything is clean and in sync
func ColorPath(s RepoStatus) string {
	switch {
	case s.Dirty:
		return color.YellowString(s.Folder)
	case s.PullNeeded || s.PushNeeded:
		return color.RedString(s.Folder)
	default:
		return color.GreenString(s.Folder)
	}
}
