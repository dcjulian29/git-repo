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

package issue

import (
	"context"

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func duplicateCmd() *cobra.Command {
	var (
		of      int
		comment string
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "duplicate <repo>#<number> --of <number>",
		Short: "Mark an issue as a duplicate of another and close it",
		Args:  cli.WithUsage(cobra.ExactArgs(1)),
		RunE: func(_ *cobra.Command, args []string) error {
			ref, err := review.ParseRef(args[0])
			if err != nil {
				return err
			}

			return review.DuplicateIssue(context.Background(), ref, of, comment, yes)
		},
	}

	cmd.Flags().IntVar(&of, "of", 0, "number of the issue this one duplicates (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "optional extra explanation")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip the confirmation prompt")

	return cmd
}
