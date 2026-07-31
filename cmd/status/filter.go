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

package status

import (
	"github.com/dcjulian29/git-repo/internal/git"
)

func filter(s git.RepoStatus) bool {
	switch {
	case dirty:
		return s.Dirty
	case pull:
		return s.PullNeeded
	case push:
		return s.PushNeeded
	case diverged:
		return s.Diverged
	case untracked:
		return s.Untracked
	case noUpstream:
		return s.NoUpstream
	case actions:
		return needsAction(s)
	default:
		return true
	}
}

// needsAction reports whether the repository requires action. Untracked files
// and a missing upstream are informational, not actionable, so they are
// excluded.
func needsAction(s git.RepoStatus) bool {
	return s.Dirty || s.PullNeeded || s.PushNeeded || s.Diverged
}
