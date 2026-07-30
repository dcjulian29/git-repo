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

// Package goversion aligns the go directive of managed repositories' go.mod
// files with the locally installed Go toolchain.
package goversion

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/fatih/color"
)

// Options controls how UpdateManaged applies the change to each repository.
type Options struct {
	// Commit records the go.mod change in a commit without pushing it.
	Commit bool

	// Push commits the change and pushes it to the remote.
	Push bool
}

// UpdateManaged rewrites the go directive of every managed repository's go.mod
// to the local Go version and removes any toolchain directive. Repositories
// without a go.mod are skipped, as are those whose go.mod or go.sum already has
// uncommitted changes.
func UpdateManaged(opts Options) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.Directory == "" {
		return errors.New(
			"configuration key 'directory' is not set; " +
				"set it with: git-repo config directory <path>")
	}

	version, err := LocalVersion()
	if err != nil {
		return err
	}

	baseDir := filesystem.ExpandHome(cfg.Directory)

	for _, repo := range cfg.Repositories {
		if repo.Managed() {
			updateRepository(git.RepositoryPath(baseDir, repo.Name), repo.Name, version, opts)
		}
	}

	return nil
}

// LocalVersion returns the version of the Go toolchain on PATH without the
// leading "go" (for example "1.26.5").
func LocalVersion() (string, error) {
	out, err := execute.ExternalProgramCapture("go", "env", "GOVERSION")
	if err != nil {
		return "", fmt.Errorf("could not determine the local Go version: %w", err)
	}

	version := normalizeVersion(out)
	if version == "" {
		return "", errors.New("the local Go version could not be determined")
	}

	return version, nil
}

// normalizeVersion trims surrounding whitespace and the leading "go" from the
// output of "go env GOVERSION".
func normalizeVersion(raw string) string {
	return strings.TrimPrefix(strings.TrimSpace(raw), "go")
}

// updateRepository applies the change to a single repository, reporting what it
// did (or why it skipped) to stdout.
func updateRepository(path, name, version string, opts Options) {
	handle := color.CyanString(name)

	if !filesystem.FileExist(filepath.Join(path, "go.mod")) {
		return
	}

	if isDirty(path) {
		fmt.Println(textformat.Warn(
			fmt.Sprintf("Skipping %s: go.mod or go.sum has uncommitted changes.", name)))

		return
	}

	current := currentVersion(path)
	if current == version {
		fmt.Printf("%s already on Go %s.\n", handle, version)

		return
	}

	if err := setVersion(path, version); err != nil {
		fmt.Println(textformat.Warn(fmt.Sprintf("Could not update %s: %v", name, err)))

		return
	}

	from := current
	if from == "" {
		from = "unset"
	}

	fmt.Printf("%s %s → %s\n", handle, from, version)

	recordChange(path, name, version, opts)
}

// recordChange commits (and optionally pushes) the go.mod change when the
// caller asked for it, limiting the commit to go.mod.
func recordChange(path, name, version string, opts Options) {
	if !opts.Commit && !opts.Push {
		return
	}

	message := "Updated to Go version " + version
	if err := git.Run(path, "commit", "-m", message, "--", "go.mod"); err != nil {
		fmt.Println(textformat.Warn(fmt.Sprintf("Could not commit %s: %v", name, err)))

		return
	}

	if !opts.Push {
		fmt.Println(color.GreenString("Committed %s.", name))

		return
	}

	if err := git.Run(path, "push"); err != nil {
		fmt.Println(textformat.Warn(fmt.Sprintf("Could not push %s: %v", name, err)))

		return
	}

	fmt.Println(color.GreenString("Committed and pushed %s.", name))
}

// isDirty reports whether go.mod or go.sum has uncommitted changes.
func isDirty(path string) bool {
	return git.CaptureOutput(path, "status", "--porcelain", "--", "go.mod", "go.sum") != ""
}

// currentVersion returns the go directive currently recorded in the
// repository's go.mod, or an empty string when it cannot be read.
func currentVersion(path string) string {
	out, err := execute.ExternalProgramCapture("go", "-C", path, "mod", "edit", "-json")
	if err != nil {
		return ""
	}

	var edit struct {
		Go string `json:"Go"`
	}

	if json.Unmarshal([]byte(out), &edit) != nil {
		return ""
	}

	return edit.Go
}

// setVersion rewrites the go directive to version and removes any toolchain
// directive, touching only go.mod.
func setVersion(path, version string) error {
	return execute.ExternalProgram("go", "-C", path, "mod", "edit",
		"-go="+version, "-toolchain=none")
}
