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

// Package initialize contains the code that creates the destination directory
// if necessary, and clones every repository that is not already present.
package initialize

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewCommand returns a cobra.Command that reads the repository list from
// the configuration file, creates the destination directory if necessary,
// and clones every repository that is not already present on disk.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Clone repositories defined in the configuration file",
		Aliases: []string{"initialize", "update"},
		Args:    cli.WithUsage(cobra.NoArgs),
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

			if len(cfg.Repositories) == 0 {
				return errors.New(
					"no repositories are configured; " +
						"add one with: git-repo config add <name> <url>",
				)
			}

			baseDir := filesystem.ExpandHome(cfg.Directory)

			if err := filesystem.EnsureDirectoryExist(baseDir); err != nil {
				return fmt.Errorf("failed to create destination directory %q: %w", baseDir, err)
			}

			toClone := planClones(baseDir, cfg.Repositories)

			if len(toClone) == 0 {
				fmt.Println(color.GreenString("\nAll configured repositories already exist. Nothing to do."))

				return nil
			}

			fmt.Printf("\nCloning %d repositor%s into %s\n\n",
				len(toClone),
				shared.Iif(len(toClone) == 1, "y", "ies"),
				color.CyanString(baseDir),
			)

			cloneErrors := cloneRepositories(baseDir, toClone)

			if len(cloneErrors) > 0 {
				return fmt.Errorf("%d clone operation(s) failed", len(cloneErrors))
			}

			fmt.Println(color.GreenString("\nDone."))

			return nil
		},
	}

	return cmd
}

// repoTarget returns the on-disk destination directory for a repository,
// treating a name containing "/" as a nested path under baseDir.
func repoTarget(baseDir, name string) string {
	target := baseDir

	if strings.Contains(name, "/") {
		for sub := range strings.SplitSeq(name, "/") {
			target = filepath.Join(target, sub)
		}
	} else {
		target = filepath.Join(target, name)
	}

	return target
}

// planClones reports which repositories still need to be cloned, printing a
// message for each one that already exists on disk.
func planClones(baseDir string, repos []config.Repository) []config.Repository {
	var toClone []config.Repository

	for _, repo := range repos {
		if filesystem.DirectoryExist(repoTarget(baseDir, repo.Name)) {
			fmt.Printf("%s  %s already exists — skipping\n",
				color.GreenString("✔"),
				color.CyanString(repo.Name))
		} else {
			toClone = append(toClone, repo)
		}
	}

	return toClone
}

// cloneRepositories clones each repository into baseDir, printing progress and
// returning the errors for any clones that failed.
func cloneRepositories(baseDir string, repos []config.Repository) []error {
	var cloneErrors []error

	for _, repo := range repos {
		if err := git.Clone(repoTarget(baseDir, repo.Name), repo.URL); err != nil {
			cloneErrors = append(cloneErrors,
				fmt.Errorf("failed to clone %s: %w", color.CyanString(repo.Name), err))
			fmt.Printf("\n%s  %s\n\n", color.RedString("✘"), err)
		} else {
			fmt.Printf("\n%s  %s\n\n", color.GreenString("✔"), repo.Name)
		}
	}

	return cloneErrors
}
