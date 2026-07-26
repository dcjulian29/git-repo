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

// Package cli holds helpers shared across the git-repo cobra commands.
package cli

import "github.com/spf13/cobra"

// WithUsage wraps a positional-argument validator so the command's usage is
// printed when validation fails. The root command sets SilenceUsage, so without
// this a bad invocation shows only a terse message such as
// "accepts 1 arg(s), received 0" with no reminder of how to call the command.
func WithUsage(validator cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := validator(cmd, args)
		if err != nil {
			_ = cmd.Usage()
		}

		return err
	}
}
