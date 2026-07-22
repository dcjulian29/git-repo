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
	"net/url"
)

// PullRequest holds the details git-repo needs to review a single pull request.
type PullRequest struct {
	// Number is the pull-request number.
	Number int

	// Title is the pull-request title.
	Title string

	// Author is the login of the user who opened it.
	Author string

	// State is the pull-request state (for example "open").
	State string

	// Draft is true when the pull request is a draft.
	Draft bool

	// HeadRef is the source branch name.
	HeadRef string

	// HeadSHA is the head commit SHA, used to look up checks.
	HeadSHA string

	// BaseRef is the target branch name.
	BaseRef string

	// Mergeable is GitHub's mergeability verdict; nil means it is still being
	// computed.
	Mergeable *bool

	// MergeableState is a coarse label such as "clean", "dirty" (conflicts),
	// "blocked", "behind", or "unstable".
	MergeableState string

	// Body is the pull-request description.
	Body string

	// URL is the html_url that opens the pull request in a browser.
	URL string
}

// HasConflicts reports whether the pull request has merge conflicts.
func (p PullRequest) HasConflicts() bool {
	return p.MergeableState == "dirty" || (p.Mergeable != nil && !*p.Mergeable)
}

type pullPayload struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	State          string `json:"state"`
	Draft          bool   `json:"draft"`
	Body           string `json:"body"`
	HTMLURL        string `json:"html_url"`
	Mergeable      *bool  `json:"mergeable"`
	MergeableState string `json:"mergeable_state"`
	User           struct {
		Login string `json:"login"`
	} `json:"user"`
	Head struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// GetPull fetches the details of a single pull request.
func (c *Client) GetPull(ctx context.Context, repo Repo, number int) (PullRequest, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	var payload pullPayload
	if err := c.getJSON(ctx, endpoint, &payload); err != nil {
		return PullRequest{}, err
	}

	return PullRequest{
		Number:         payload.Number,
		Title:          payload.Title,
		Author:         payload.User.Login,
		State:          payload.State,
		Draft:          payload.Draft,
		HeadRef:        payload.Head.Ref,
		HeadSHA:        payload.Head.SHA,
		BaseRef:        payload.Base.Ref,
		Mergeable:      payload.Mergeable,
		MergeableState: payload.MergeableState,
		Body:           payload.Body,
		URL:            payload.HTMLURL,
	}, nil
}
