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

// Package cmd contains all Cobra sub-commands for the git-repo CLI.
// The root command wires together all sub-commands so that they are available
// from the single git-repo binary.
package cmd

import (
	"github.com/dcjulian29/git-repo/cmd/config"
	"github.com/dcjulian29/git-repo/cmd/gocmd"
	"github.com/dcjulian29/git-repo/cmd/initialize"
	"github.com/dcjulian29/git-repo/cmd/issue"
	"github.com/dcjulian29/git-repo/cmd/label"
	"github.com/dcjulian29/git-repo/cmd/pr"
	"github.com/dcjulian29/git-repo/cmd/status"
	"github.com/dcjulian29/git-repo/cmd/synchronize"
	"github.com/spf13/cobra"
	"go.szostok.io/version/extension"
)

var rootCmd = &cobra.Command{
	Use:   "git-repo",
	Short: "Manage multiple Git repositories",
	Long: `git-repo is a CLI tool that lets you inspect the
status of, synchronize, initialize, and configure multiple
local Git repositories declared in ~/.config/git-repo.yml.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.AddCommand(config.NewCommand())
	rootCmd.AddCommand(gocmd.NewCommand())
	rootCmd.AddCommand(initialize.NewCommand())
	rootCmd.AddCommand(issue.NewCommand())
	rootCmd.AddCommand(label.NewCommand())
	rootCmd.AddCommand(pr.NewCommand())
	rootCmd.AddCommand(status.NewCommand())
	rootCmd.AddCommand(synchronize.NewCommand())
}

// Execute is the entry-point called by main. It builds the command tree and
// runs the appropriate sub-command based on os.Args. On failure it returns the
// error so that the caller can report it and set the process exit code.
func Execute() error {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "git-repo"),
		),
	)

	return rootCmd.Execute()
}
