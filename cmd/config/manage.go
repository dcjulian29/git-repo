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

	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func manageCmd() *cobra.Command {
	return manageStateCmd(
		"manage <name>",
		"Mark a configured repository as managed by the pr and issue commands",
		true,
	)
}

func unmanageCmd() *cobra.Command {
	return manageStateCmd(
		"unmanage <name>",
		"Mark a configured repository so the pr and issue commands skip it",
		false,
	)
}

// manageStateCmd builds a command that sets the managed state of a single
// configured repository identified by name.
func manageStateCmd(use, short string, managed bool) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				if err := cmd.Usage(); err != nil {
					return err
				}

				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return setManaged(args[0], managed)
		},
	}
}

// setManaged updates the managed flag of the named repository and saves the
// configuration. Managed repositories clear the flag (returning to the
// default); unmanaged repositories set it to false.
func setManaged(name string, managed bool) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	for i := range cfg.Repositories {
		if !strings.EqualFold(cfg.Repositories[i].Name, name) {
			continue
		}

		if managed {
			cfg.Repositories[i].Manage = nil
		} else {
			value := false
			cfg.Repositories[i].Manage = &value
		}

		if err := config.Save(&cfg); err != nil {
			return err
		}

		state := "managed"
		if !managed {
			state = "unmanaged"
		}

		fmt.Printf("Marked %s as %s.\n", color.CyanString(name), state)

		return nil
	}

	return fmt.Errorf("repository %s not found in the configuration", color.CyanString(name))
}
