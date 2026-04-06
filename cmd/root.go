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
// The root command initializes Viper configuration and wires together all
// sub-commands so that they are available from the single git-repo binary.
package cmd

import (
	"fmt"
	"os"

	"github.com/dcjulian29/git-repo/cmd/config"
	"github.com/dcjulian29/git-repo/cmd/initialize"
	"github.com/dcjulian29/go-toolbox/textformat"
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

// Execute is the entry-point called by main. It builds the command tree and
// runs the appropriate sub-command based on os.Args.
func Execute() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "git-repo"),
		),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "\n"+textformat.Fatal(err.Error()))
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(config.NewCommand())
	rootCmd.AddCommand(initialize.NewCommand())
}
