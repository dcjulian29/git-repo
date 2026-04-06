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
	"os"

	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured repositories",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := shared.GetSettings()
			if err != nil {
				return err
			}

			return showRepositories(&cfg)
		},
	}

	return cmd
}

func showRepositories(cfg *shared.Config) error {
	if len(cfg.Repositories) == 0 {
		fmt.Println(textformat.Warn("No repositories are configured."))

		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "URL"})

	for _, r := range cfg.Repositories {
		if err := table.Append([]string{color.GreenString(r.Name), color.CyanString(r.URL)}); err != nil {
			return err
		}
	}

	return table.Render()
}
