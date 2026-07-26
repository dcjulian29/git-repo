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

package pr

import (
	"context"

	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var (
		open       bool
		closed     bool
		merged     bool
		draft      bool
		mergeState string
		asJSON     bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull requests across managed repositories",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			mode := "open"

			switch {
			case draft:
				mode = "draft"
			case closed:
				mode = "closed"
			case merged:
				mode = "merged"
			}

			return review.ListPulls(context.Background(), mode, mergeState, asJSON)
		},
	}

	cmd.Flags().BoolVar(&open, "open", false, "list open pull requests (default)")
	cmd.Flags().BoolVar(&closed, "closed", false, "list closed (unmerged) pull requests")
	cmd.Flags().BoolVar(&merged, "merged", false, "list merged pull requests")
	cmd.Flags().BoolVar(&draft, "draft", false, "list draft pull requests")
	cmd.Flags().StringVar(&mergeState, "merge-state", "",
		"only pull requests with this merge state (e.g. clean, unstable, unknown; default all)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "output as JSON")
	cmd.MarkFlagsMutuallyExclusive("open", "closed", "merged", "draft")

	return cmd
}
