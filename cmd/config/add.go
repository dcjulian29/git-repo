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

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add a repository to the configuration",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			name, url := args[0], args[1]

			cfg, err := shared.GetSettings()
			if err != nil {
				return err
			}

			for _, r := range cfg.Repositories {
				if strings.EqualFold(r.Name, name) {
					return fmt.Errorf("repository %q already exists in the configuration", name)
				}
			}

			cfg.Repositories = append(cfg.Repositories, shared.Repository{Name: name, URL: url})

			if err := shared.SaveSettings(&cfg); err != nil {
				return err
			}

			fmt.Printf("Added %s → %s\n",
				color.GreenString(name),
				color.CyanString(url),
			)

			return nil
		},
	}

	return cmd
}
