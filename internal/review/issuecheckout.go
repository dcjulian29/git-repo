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

	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/fatih/color"
)

// CheckoutIssue creates a working branch for the issue in the local clone and on
// origin, opens a draft pull request linked to the issue via "Closes #<n>", and
// leaves the working tree on the new branch unless returnToBase is set. When
// branch is empty it defaults to "issue/<number>".
func CheckoutIssue(ctx context.Context, ref Ref, branch string, returnToBase bool) error {
	if !filesystem.DirectoryExist(ref.Path) {
		return fmt.Errorf("repository %s is not cloned at %s; run 'git-repo init' first", ref.Name, ref.Path)
	}

	if branch == "" {
		branch = fmt.Sprintf("issue/%d", ref.Number)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	base, err := client.DefaultBranch(ctx, ref.Repo)
	if err != nil {
		return err
	}

	steps := [][]string{
		{"fetch", "origin"},
		{"checkout", "-b", branch, "origin/" + base},
		{"commit", "--allow-empty", "-m", fmt.Sprintf("Start work on #%d: %s", ref.Number, issue.Title)},
		{"push", "-u", "origin", branch},
	}

	for _, args := range steps {
		if runErr := git.Run(ref.Path, args...); runErr != nil {
			return runErr
		}
	}

	pull, err := client.CreatePull(ctx, ref.Repo, github.CreatePullParams{
		Title: issue.Title,
		Head:  branch,
		Base:  base,
		Body:  fmt.Sprintf("Closes #%d", ref.Number),
		Draft: true,
	})
	if err != nil {
		return err
	}

	fmt.Println(color.GreenString("Opened draft PR %s#%d for issue #%d — %s",
		ref.Name, pull.Number, ref.Number, pull.URL))

	if returnToBase {
		return git.Run(ref.Path, "checkout", base)
	}

	return nil
}
