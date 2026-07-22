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

package github

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
)

// Token resolves the GitHub API token. It prefers the GITHUB_TOKEN environment
// variable and otherwise falls back to running "gh auth token". The token is
// deliberately never read from the git-repo configuration file, which may be
// shared publicly.
func Token() (string, error) {
	if t := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); t != "" {
		return t, nil
	}

	out, err := execute.ExternalProgramCapture("gh", "auth", "token")
	if err != nil {
		return "", fmt.Errorf(
			"no GitHub token found; set GITHUB_TOKEN or authenticate with 'gh auth login': %w", err)
	}

	token := strings.TrimSpace(out)
	if token == "" {
		return "", errors.New("'gh auth token' returned an empty token; run 'gh auth login'")
	}

	return token, nil
}
