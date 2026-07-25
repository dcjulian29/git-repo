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

// Package issue contains the sub-commands for working with open issues across
// the managed repositories.
package issue

import "github.com/spf13/cobra"

// NewCommand returns the "issue" command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "Work with open issues across managed repositories",
		Aliases: []string{"issues"},
	}

	cmd.AddCommand(acceptCmd())
	cmd.AddCommand(checkoutCmd())
	cmd.AddCommand(closeCmd())
	cmd.AddCommand(duplicateCmd())
	cmd.AddCommand(labelCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(openCmd())
	cmd.AddCommand(rejectCmd())
	cmd.AddCommand(showCmd())
	cmd.AddCommand(triageCmd())

	return cmd
}
