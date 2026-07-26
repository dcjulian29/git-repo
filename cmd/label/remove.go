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

package label

import (
	"context"

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:     "remove <repo> <name>",
		Short:   "Delete a label from a repository",
		Aliases: []string{"rm", "delete", "del"},
		Args:    cli.WithUsage(cobra.ExactArgs(2)),
		RunE: func(_ *cobra.Command, args []string) error {
			return review.RemoveLabel(context.Background(), args[0], args[1], yes)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip the confirmation prompt")

	return cmd
}
