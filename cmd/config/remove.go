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

package config

import (
	"fmt"
	"strings"

	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <name>",
		Short:   "Remove a repository from the configuration by name",
		Aliases: []string{"rm", "delete", "del"},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := shared.GetSettings()
			if err != nil {
				return err
			}

			original := len(cfg.Repositories)
			var kept []shared.Repository

			for _, r := range cfg.Repositories {
				if !strings.EqualFold(r.Name, name) {
					kept = append(kept, r)
				}
			}

			if len(kept) == original {
				return fmt.Errorf("repository %s not found in the configuration", color.CyanString(name))
			}

			cfg.Repositories = kept

			if err := shared.SaveSettings(&cfg); err != nil {
				return err
			}

			fmt.Printf("Removed %s from the configuration.\n", color.CyanString(name))

			return nil
		},
	}

	return cmd
}
