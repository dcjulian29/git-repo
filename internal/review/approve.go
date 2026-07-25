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

package review

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/fatih/color"
)

// ApprovePull shows the pull-request summary, then (unless skipConfirm is set)
// asks for confirmation before merging it with the given method. Pull requests
// with conflicts are refused, since they cannot be merged as-is.
func ApprovePull(ctx context.Context, ref Ref, method string, skipConfirm bool) error {
	if !slices.Contains(github.MergeMethods, method) {
		return fmt.Errorf("invalid merge method %q; expected one of %v", method, github.MergeMethods)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	pull, err := describe(ctx, client, ref)
	if err != nil {
		return err
	}

	if pull.HasConflicts() {
		return fmt.Errorf(
			"%s#%d has conflicts and cannot be merged; run 'git-repo pr checkout %s#%d' to resolve them",
			ref.Name, ref.Number, ref.Name, ref.Number)
	}

	if !skipConfirm {
		prompt := fmt.Sprintf("\nMerge %s using %s? [y/N] ",
			color.CyanString("%s#%d", ref.Name, ref.Number), method)

		confirmed, err := shared.Confirm(os.Stdin, os.Stdout, prompt)
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Aborted.")

			return nil
		}
	}

	if err := client.MergePull(ctx, ref.Repo, ref.Number, method); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Merged %s#%d.", ref.Name, ref.Number))

	return nil
}
