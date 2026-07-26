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
	"fmt"

	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/spf13/cobra"
)

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open <repo>#<number>",
		Short: "Open an issue in the default browser",
		Args:  cli.WithUsage(cobra.ExactArgs(1)),
		RunE: func(_ *cobra.Command, args []string) error {
			ref, err := review.ParseRef(args[0])
			if err != nil {
				return err
			}

			target := fmt.Sprintf("https://github.com/%s/%s/issues/%d",
				ref.Repo.Owner, ref.Repo.Name, ref.Number)

			fmt.Printf("Opening %s\n", target)

			return shared.OpenBrowser(target)
		},
	}
}
