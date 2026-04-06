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

// Package config groups show, directory, add, remove, and list sub-commands under the
// single "config" namespace so that all CRUD operations are co-located.
package config

import (
	"github.com/spf13/cobra"
)

// NewCommand returns a cobra.Command that provides configuration management.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and manage the git-repo configuration file",
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(directoryCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(removeCmd())
	cmd.AddCommand(showCmd())

	return cmd
}
