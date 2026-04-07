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

// RepoStatus holds the computed status of a single Git repository.
// All boolean fields are derived from porcelain git output so that the
// status and sync commands share a single data-gathering path.
type RepoStatus struct {
	// Folder is the absolute path to the repository root.
	Folder string

	// Dirty is true when "git diff --stat" reports uncommitted changes.
	Dirty bool

	// PushNeeded is true when the local branch is ahead of the remote.
	PushNeeded bool

	// PullNeeded is true when the remote branch is ahead of the local HEAD.
	PullNeeded bool

	// Diverged is true when local and remote have diverged from each other.
	Diverged bool

	// Untracked is true when there are untracked files in the working tree.
	Untracked bool
}
