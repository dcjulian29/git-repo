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
	"fmt"
	"strings"
)

// Repo identifies a GitHub repository by its owner and name.
type Repo struct {
	// Owner is the user or organization that owns the repository.
	Owner string

	// Name is the repository name, without any ".git" suffix.
	Name string
}

// String returns the canonical "owner/name" slug for the repository.
func (r Repo) String() string {
	return r.Owner + "/" + r.Name
}

// prefixes maps the accepted github.com remote URL prefixes to strip before
// parsing the owner and name.
var githubPrefixes = []string{
	"git@github.com:",
	"ssh://git@github.com/",
	"https://github.com/",
	"http://github.com/",
}

// ParseRepo extracts the owner and name from a GitHub remote URL. It accepts
// both HTTPS ("https://github.com/owner/name.git") and SSH
// ("git@github.com:owner/name.git") forms. It returns an error for URLs that
// are not hosted on github.com, so callers can skip third-party hosts.
func ParseRepo(remote string) (Repo, error) {
	trimmed := strings.TrimSpace(remote)

	path := ""
	matched := false

	for _, prefix := range githubPrefixes {
		if rest, found := strings.CutPrefix(trimmed, prefix); found {
			path = rest
			matched = true

			break
		}
	}

	if !matched {
		return Repo{}, fmt.Errorf("not a github.com repository: %q", remote)
	}

	path = strings.TrimSuffix(path, ".git")
	path = strings.Trim(path, "/")

	owner, name, found := strings.Cut(path, "/")
	if !found || owner == "" || name == "" {
		return Repo{}, fmt.Errorf("could not parse owner and name from: %q", remote)
	}

	if strings.Contains(name, "/") {
		return Repo{}, fmt.Errorf("unexpected extra path in github repository url: %q", remote)
	}

	return Repo{Owner: owner, Name: name}, nil
}
