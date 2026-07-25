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

	"github.com/dcjulian29/git-repo/internal/review"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <repo>",
		Short: "List the labels defined in a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return review.ListLabels(context.Background(), args[0])
		},
	}
}
