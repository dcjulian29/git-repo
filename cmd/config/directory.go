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

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

func directoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "directory <path>",
		Aliases: []string{"set-dir"},
		Short:   "Set or update the managed directory",
		Args:    cli.WithUsage(cobra.ExactArgs(1)),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			cfg.Directory = args[0]

			if err := config.Save(&cfg); err != nil {
				return err
			}

			fmt.Printf("Directory set to %s\n", textformat.Teal(args[0]))

			return nil
		},
	}

	return cmd
}
