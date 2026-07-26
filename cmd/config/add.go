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

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	var noManage bool

	cmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add a repository to the configuration",
		Args:  cli.WithUsage(cobra.ExactArgs(2)),
		RunE: func(_ *cobra.Command, args []string) error {
			name, url := args[0], args[1]

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			for _, r := range cfg.Repositories {
				if strings.EqualFold(r.Name, name) {
					return fmt.Errorf("repository %q already exists in the configuration", name)
				}
			}

			repo := config.Repository{Name: name, URL: url}
			if noManage {
				value := false
				repo.Manage = &value
			}

			cfg.Repositories = append(cfg.Repositories, repo)

			if err := config.Save(&cfg); err != nil {
				return err
			}

			fmt.Printf("Added %s → %s\n",
				color.GreenString(name),
				color.CyanString(url),
			)

			return nil
		},
	}

	cmd.Flags().BoolVar(&noManage, "no-manage", false,
		"add the repository as unmanaged so pr and issue commands skip it")

	return cmd
}
