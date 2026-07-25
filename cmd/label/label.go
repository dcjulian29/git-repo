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

// Package label contains the sub-commands for managing the labels of a
// configured repository.
package label

import "github.com/spf13/cobra"

// NewCommand returns the "label" command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "label",
		Short:   "Manage the labels of a configured repository",
		Aliases: []string{"labels"},
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(removeCmd())

	return cmd
}
