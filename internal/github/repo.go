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
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DefaultBranch returns the repository's default branch name.
func (c *Client) DefaultBranch(ctx context.Context, repo Repo) (string, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name))

	var payload struct {
		DefaultBranch string `json:"default_branch"`
	}

	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &payload); err != nil {
		return "", err
	}

	return payload.DefaultBranch, nil
}

// CreatePullParams describes a pull request to open.
type CreatePullParams struct {
	// Title is the pull-request title.
	Title string

	// Head is the source branch.
	Head string

	// Base is the target branch.
	Base string

	// Body is the pull-request description.
	Body string

	// Draft opens the pull request as a draft when true.
	Draft bool
}

// CreatePull opens a new pull request and returns it.
func (c *Client) CreatePull(ctx context.Context, repo Repo, params CreatePullParams) (PullRequest, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name))

	body := map[string]any{
		"title": params.Title,
		"head":  params.Head,
		"base":  params.Base,
		"body":  params.Body,
		"draft": params.Draft,
	}

	var payload pullPayload
	if err := c.doJSON(ctx, http.MethodPost, endpoint, body, &payload); err != nil {
		return PullRequest{}, err
	}

	return payload.toPullRequest(), nil
}
