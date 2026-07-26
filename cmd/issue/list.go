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
	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var (
		open     bool
		closed   bool
		all      bool
		labels   []string
		assignee string
		asJSON   bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues across managed repositories",
		Args:  cli.WithUsage(cobra.NoArgs),
		RunE: func(_ *cobra.Command, _ []string) error {
			state := "open"

			switch {
			case closed:
				state = "closed"
			case all:
				state = "all"
			default:
				// no flag set: keep the default "open" state
			}

			opts := github.ListOptions{State: state, Labels: labels, Assignee: assignee}

			return review.ListIssues(context.Background(), opts, asJSON)
		},
	}

	cmd.Flags().BoolVar(&open, "open", false, "list open issues (default)")
	cmd.Flags().BoolVar(&closed, "closed", false, "list closed issues")
	cmd.Flags().BoolVar(&all, "all", false, "list open and closed issues")
	cmd.Flags().StringSliceVar(&labels, "label", nil, "only issues carrying the given label (repeatable)")
	cmd.Flags().StringVar(&assignee, "assignee", "", `only issues assigned to this login ("@me" for yourself)`)
	cmd.Flags().BoolVar(&asJSON, "json", false, "output as JSON")
	cmd.MarkFlagsMutuallyExclusive("open", "closed", "all")

	return cmd
}
