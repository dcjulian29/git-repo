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

// Package synchronize synchronizes every discovered Git repository (fetch → pull --rebase → push).
package synchronize

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewCommand returns a cobra.Command that walks a directory tree and
// synchronizes every discovered Git repository (fetch → pull --rebase → push).
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "synchronize",
		Short:   "Synchronize all managed git repositories",
		Aliases: []string{"sync", "synchronise"},
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

			results := shared.ParallelMap(dirs, shared.DefaultConcurrency(), git.Synchronize)

			_ = spinner.Stop()

			sort.Slice(results, func(i, j int) bool {
				return strings.Compare(results[i].Folder, results[j].Folder) < 0
			})

			for _, r := range results {
				if len(strings.TrimSpace(r.Output)) > 0 {
					fmt.Printf("\n%s\n%s\n", color.CyanString(r.Folder), r.Output)
				}
			}

			return nil
		},
	}

	return cmd
}
