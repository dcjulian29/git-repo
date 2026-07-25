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

	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func checkoutCmd() *cobra.Command {
	var (
		branch   string
		noSwitch bool
	)

	cmd := &cobra.Command{
		Use:   "checkout <repo>#<number>",
		Short: "Create a branch and draft PR to start work on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ref, err := review.ParseRef(args[0])
			if err != nil {
				return err
			}

			return review.CheckoutIssue(context.Background(), ref, branch, noSwitch)
		},
	}

	cmd.Flags().StringVar(&branch, "branch", "", "branch name to create (defaults to issue/<number>)")
	cmd.Flags().BoolVar(&noSwitch, "no-switch", false, "leave the working tree on its current branch")

	return cmd
}
