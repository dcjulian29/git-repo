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

// Package status determines the current status of every discovered Git repository.
package status

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

var (
	actions    bool
	dirty      bool
	push       bool
	pull       bool
	diverged   bool
	untracked  bool
	noUpstream bool
)

// NewCommand returns a cobra.Command that walks a directory tree and prints a
// color-coded table showing the status of every discovered Git repository.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of all managed git repositories",
		Args:  cli.WithUsage(cobra.NoArgs),
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if cfg.Directory == "" {
				return errors.New(
					"configuration key 'directory' is not set; " +
						"set it with: git-repo config directory <path>",
				)
			}

			root := filesystem.ExpandHome(cfg.Directory)

			results, err := collectStatuses(root)
			if err != nil {
				return err
			}

			slices.SortFunc(results, func(a, b git.RepoStatus) int {
				return cmp.Compare(a.Folder, b.Folder)
			})

			render(results)

			return nil
		},
	}

	cmd.Flags().BoolVar(&actions, "actions", false, "show only repositories that require action")
	cmd.Flags().BoolVar(&dirty, "dirty", false, "show only dirty repositories")
	cmd.Flags().BoolVar(&push, "push", false, "show only repositories that need a push")
	cmd.Flags().BoolVar(&pull, "pull", false, "show only repositories that need a pull")
	cmd.Flags().BoolVar(&diverged, "diverged", false, "show only diverged repositories")
	cmd.Flags().BoolVar(&untracked, "untracked", false, "show only repositories with untracked files")
	cmd.Flags().BoolVar(&noUpstream, "no-upstream", false, "show only repositories with no upstream tracking branch")

	return cmd
}

// collectStatuses walks root for Git repositories, computes their status in
// parallel while a spinner is shown, and returns the statuses that pass the
// active filters.
func collectStatuses(root string) ([]git.RepoStatus, error) {
	spinner, err := shared.NewSpinner()
	if err != nil {
		return nil, fmt.Errorf("failed to create spinner: %w", err)
	}

	_ = spinner.Start()

	dirs, err := git.FindGitRepositories(root)
	if err != nil {
		_ = spinner.Stop()

		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	statuses := shared.ParallelMap(dirs, shared.DefaultConcurrency(), git.ComputeStatus)

	_ = spinner.Stop()

	var results []git.RepoStatus

	for _, s := range statuses {
		if filter(s) {
			results = append(results, s)
		}
	}

	return results, nil
}
