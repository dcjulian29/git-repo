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
	"strings"

	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// CheckoutPull fetches the pull request's branch into the local clone so it can
// be inspected or fixed. Same-repository pull requests become a tracking branch;
// pull requests from forks are fetched into a local "pr/<number>" branch.
func CheckoutPull(ctx context.Context, ref Ref) error {
	if !filesystem.DirectoryExist(ref.Path) {
		return fmt.Errorf("repository %s is not cloned at %s; run 'git-repo init' first", ref.Name, ref.Path)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	pull, err := client.GetPull(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	for _, args := range checkoutPlan(ref, pull) {
		if err := git.Run(ref.Path, args...); err != nil {
			return err
		}
	}

	return nil
}

// checkoutPlan returns the sequence of git commands used to check out the pull
// request locally. Same-repository pull requests are checked out as a tracking
// branch (so fixes can be pushed back); fork pull requests are fetched into a
// local pr/<number> branch.
func checkoutPlan(ref Ref, pull github.PullRequest) [][]string {
	if strings.EqualFold(pull.HeadRepo, ref.Repo.String()) && pull.HeadRef != "" {
		return [][]string{
			{"fetch", "origin"},
			{"checkout", pull.HeadRef},
		}
	}

	local := fmt.Sprintf("pr/%d", ref.Number)
	spec := fmt.Sprintf("pull/%d/head:%s", ref.Number, local)

	return [][]string{
		{"fetch", "origin", spec},
		{"checkout", local},
	}
}
