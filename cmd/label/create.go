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

func createCmd() *cobra.Command {
	var (
		labelColor  string
		description string
	)

	cmd := &cobra.Command{
		Use:   "create <repo> <name>",
		Short: "Create a label in a repository",
		Args:  cli.WithUsage(cobra.ExactArgs(2)),
		RunE: func(_ *cobra.Command, args []string) error {
			return review.CreateLabel(context.Background(), args[0], args[1], labelColor, description)
		},
	}

	cmd.Flags().StringVar(&labelColor, "color", "", "six-character hex color (defaults to a neutral gray)")
	cmd.Flags().StringVar(&description, "description", "", "label description")

	return cmd
}
