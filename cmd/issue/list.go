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
	"os"

	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List open issues across managed repositories",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			report, err := review.Collect(context.Background())
			if err != nil {
				return err
			}

			review.PrintWarnings(os.Stderr, report)

			issues := report.Issues()

			if asJSON {
				return review.RenderJSON(issues)
			}

			return review.RenderTable(issues, "ISSUE", "No open issues found.")
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "output as JSON")

	return cmd
}
