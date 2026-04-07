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
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

var (
	actions   bool
	dirty     bool
	push      bool
	pull      bool
	diverged  bool
	untracked bool
)

// NewCommand returns a cobra.Command that walks a directory tree and prints a
// colour-coded table showing the status of every discovered Git repository.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of all managed git repositories",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.NoArgs(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := shared.GetSettings()
			if err != nil {
				return err
			}

			root := filesystem.ExpandHome(cfg.Directory)

			spinner, err := shared.NewSpinner()
			if err != nil {
				return fmt.Errorf("failed to create spinner: %w", err)
			}

			_ = spinner.Start()

			dirs, err := git.FindGitRepositories(root)
			if err != nil {
				_ = spinner.Stop()
				return fmt.Errorf("failed to walk directory: %w", err)
			}

			var (
				wg      sync.WaitGroup
				mu      sync.Mutex
				results []git.RepoStatus
				ch      = make(chan git.RepoStatus, len(dirs))
			)

			for _, d := range dirs {
				wg.Add(1)

				go func(dir string) {
					defer wg.Done()
					ch <- git.ComputeStatus(dir)
				}(d)
			}

			go func() {
				wg.Wait()
				close(ch)
			}()

			for s := range ch {
				if filter(s) {
					mu.Lock()
					results = append(results, s)
					mu.Unlock()
				}
			}

			_ = spinner.Stop()

			sort.Slice(results, func(i, j int) bool {
				return strings.Compare(results[i].Folder, results[j].Folder) < 0
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

	return cmd
}
