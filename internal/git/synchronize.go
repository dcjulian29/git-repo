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

import (
	"fmt"

	"github.com/dcjulian29/go-toolbox/textformat"
)

// Synchronize performs a full sync cycle on a single repository: fetch, then
// pull (rebase + prune + submodules) if behind, then push if ahead. When the
// working tree is dirty the sync is skipped with a notice in Output.
func Synchronize(path string) SyncResult {
	_ = CaptureOutput(path, "fetch")

	dirty := CaptureOutput(path, "diff", "--stat")
	local := CaptureOutput(path, "rev-parse", "@")
	remote := CaptureOutput(path, "rev-parse", "@{u}")
	base := CaptureOutput(path, "merge-base", "@", "@{u}")

	pull, push := false, false

	if local != remote {
		switch {
		case local == base:
			pull = true
		case remote == base:
			push = true
		default:
			pull = true
			push = true
		}
	}

	var output string

	if len(dirty) > 0 {
		output = textformat.Info("Skipping — working tree is not clean.")
	} else {
		if pull {
			output = CaptureError(path, "pull", "--rebase", "--prune", "--recurse-submodules=yes")
		}

		if push {
			pushOut := CaptureError(path, "push")
			if len(output) > 0 {
				output = fmt.Sprintf("%s\n\n%s", output, pushOut)
			} else {
				output = pushOut
			}
		}
	}

	return SyncResult{Folder: path, Output: output}
}
