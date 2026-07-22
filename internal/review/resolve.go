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

package review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/github"
)

// Ref identifies a single pull request or issue within a configured repository.
type Ref struct {
	// Name is the configuration name of the repository.
	Name string

	// Repo is the parsed GitHub owner and name.
	Repo github.Repo

	// Number is the pull-request or issue number.
	Number int
}

// ParseRef resolves a "<repo>#<number>" handle to a configured github.com
// repository and item number. The repository must exist in the configuration.
func ParseRef(handle string) (Ref, error) {
	name, digits, found := strings.Cut(handle, "#")
	if !found {
		return Ref{}, fmt.Errorf("invalid reference %q; expected <repo>#<number>", handle)
	}

	number, err := strconv.Atoi(strings.TrimSpace(digits))
	if err != nil || number <= 0 {
		return Ref{}, fmt.Errorf("invalid number in reference %q; expected <repo>#<number>", handle)
	}

	name = strings.TrimSpace(name)

	cfg, err := config.Load()
	if err != nil {
		return Ref{}, err
	}

	for _, repository := range cfg.Repositories {
		if !strings.EqualFold(repository.Name, name) {
			continue
		}

		repo, err := github.ParseRepo(repository.URL)
		if err != nil {
			return Ref{}, fmt.Errorf("%s is not a github.com repository", repository.Name)
		}

		return Ref{Name: repository.Name, Repo: repo, Number: number}, nil
	}

	return Ref{}, fmt.Errorf("repository %q is not in the configuration", name)
}
