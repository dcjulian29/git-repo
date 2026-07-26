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
	"cmp"
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/dcjulian29/git-repo/internal/config"
	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
)

type (
	// NamedItem is an open item together with the display name of the repository
	// it belongs to.
	NamedItem struct {
		// Repo is the configuration name of the owning repository.
		Repo string `json:"repo"`

		// ghRepo is the parsed repository, used to enrich pull requests with data
		// the listing endpoint does not return. It is not serialised.
		ghRepo github.Repo

		github.Item
	}

	// Report is the outcome of collecting open items across managed repositories.
	Report struct {
		// Results holds one entry per fetched repository, including any error.
		Results []github.Result

		// Skipped names the managed repositories that were not on github.com.
		Skipped []string
	}
)

// Collect gathers the issues and pull requests matching opts for every managed
// github.com repository, using the given client. Managed repositories whose URL
// is not on github.com are recorded in Report.Skipped instead of being fetched.
// An assignee of "@me" is resolved to the authenticated user's login.
func Collect(ctx context.Context, client *github.Client, opts github.ListOptions) (Report, error) {
	cfg, err := config.Load()
	if err != nil {
		return Report{}, err
	}

	var (
		targets []github.Target
		skipped []string
	)

	for _, repository := range cfg.Repositories {
		if !repository.Managed() {
			continue
		}

		repo, err := github.ParseRepo(repository.URL)
		if err != nil {
			skipped = append(skipped, repository.Name)

			continue
		}

		targets = append(targets, github.Target{Name: repository.Name, Repo: repo})
	}

	if len(targets) == 0 {
		return Report{Skipped: skipped}, errors.New("no managed github.com repositories are configured")
	}

	if opts.Assignee == "@me" {
		login, err := client.Viewer(ctx)
		if err != nil {
			return Report{}, err
		}

		opts.Assignee = login
	}

	results := client.Gather(ctx, targets, shared.DefaultConcurrency(), opts)

	return Report{Results: results, Skipped: skipped}, nil
}

// PullRequests returns every open pull request across the report, sorted by
// repository name and then by number.
func (r Report) PullRequests() []NamedItem {
	return r.filter(true)
}

// Issues returns every open issue across the report, sorted by repository name
// and then by number.
func (r Report) Issues() []NamedItem {
	return r.filter(false)
}

// Failures returns the results whose repositories could not be fetched.
func (r Report) Failures() []github.Result {
	var failed []github.Result

	for _, res := range r.Results {
		if res.Err != nil {
			failed = append(failed, res)
		}
	}

	return failed
}

// filter flattens the successful results into named items of the requested
// kind, sorted by repository name and then by number.
func (r Report) filter(pull bool) []NamedItem {
	var items []NamedItem

	for _, res := range r.Results {
		if res.Err != nil {
			continue
		}

		for _, item := range res.Items {
			if item.IsPull == pull {
				items = append(items, NamedItem{Repo: res.Target.Name, Item: item, ghRepo: res.Target.Repo})
			}
		}
	}

	slices.SortStableFunc(items, func(a, b NamedItem) int {
		return cmp.Or(
			cmp.Compare(strings.ToLower(a.Repo), strings.ToLower(b.Repo)),
			cmp.Compare(a.Number, b.Number),
		)
	})

	return items
}
