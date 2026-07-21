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
// a single repository rooted at path. When the current branch has no upstream
// tracking branch the push/pull/diverged comparison is skipped and NoUpstream
// is set instead, because there is no remote to compare against.
func ComputeStatus(path string) RepoStatus {
	_ = CaptureOutput(path, "fetch")

	dirty := CaptureOutput(path, "diff", "--stat")
	untracked := CaptureOutput(path, "ls-files", "--others", "--exclude-standard")

	status := RepoStatus{
		Folder:    path,
		Dirty:     len(dirty) > 0,
		Untracked: len(untracked) > 0,
	}

	if !HasUpstream(path) {
		status.NoUpstream = true

		return status
	}

	local := CaptureOutput(path, "rev-parse", "@")
	remote := CaptureOutput(path, "rev-parse", "@{u}")
	base := CaptureOutput(path, "merge-base", "@", "@{u}")

	if local != remote {
		switch {
		case local == base:
			status.PullNeeded = true
		case remote == base:
			status.PushNeeded = true
		default:
			status.Diverged = true
		}
	}

	return status
}
