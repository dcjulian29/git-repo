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

// ComputeStatus fetches from the remote and derives the full status for
// a single repository rooted at path.
func ComputeStatus(path string) RepoStatus {
	_ = CaptureOutput(path, "fetch")

	dirty := CaptureOutput(path, "diff", "--stat")
	untracked := CaptureOutput(path, "ls-files", "--others", "--exclude-standard")
	local := CaptureOutput(path, "rev-parse", "@")
	remote := CaptureOutput(path, "rev-parse", "@{u}")
	base := CaptureOutput(path, "merge-base", "@", "@{u}")

	pull, push, diverged := false, false, false

	if local != remote {
		switch {
		case local == base:
			pull = true
		case remote == base:
			push = true
		default:
			diverged = true
		}
	}

	return RepoStatus{
		Folder:     path,
		Dirty:      len(dirty) > 0,
		PushNeeded: push,
		PullNeeded: pull,
		Diverged:   diverged,
		Untracked:  len(untracked) > 0,
	}
}
