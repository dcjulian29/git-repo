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

package gocmd

import (
	"github.com/dcjulian29/git-repo/internal/cli"
	"github.com/dcjulian29/git-repo/internal/goversion"
	"github.com/spf13/cobra"
)

// updateCmd builds the "go update" command that sets each managed repository's
// go.mod Go version to the locally installed toolchain.
func updateCmd() *cobra.Command {
	var opts goversion.Options

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Set each managed repository's go.mod Go version to the local toolchain",
		Long: `Rewrite the go directive of every managed repository's go.mod to the
locally installed Go version and remove any toolchain directive, so that local
and CI builds use the same Go. Only go.mod is changed; other dependencies are
left untouched. A repository is skipped when it has no go.mod or when its go.mod
or go.sum already has uncommitted changes.

By default the change is left in the working tree for review. Use --commit to
commit it (without pushing), or --push to commit and push.`,
		Args: cli.WithUsage(cobra.NoArgs),
		RunE: func(_ *cobra.Command, _ []string) error {
			return goversion.UpdateManaged(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Commit, "commit", false, "commit the change without pushing")
	cmd.Flags().BoolVar(&opts.Push, "push", false, "commit and push the change")

	return cmd
}
