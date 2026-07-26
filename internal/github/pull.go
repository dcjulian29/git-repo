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

// PullRequest holds the details git-repo needs to review a single pull request.
type PullRequest struct {
	// Number is the pull-request number.
	Number int

	// NodeID is the GraphQL global node ID, needed for GraphQL-only mutations
	// such as marking a draft ready for review.
	NodeID string

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

	// HeadRepo is the "owner/name" of the repository the head branch lives in.
	// It differs from the base repository when the pull request comes from a
	// fork.
	HeadRepo string
}

// HasConflicts reports whether the pull request has merge conflicts.
func (p PullRequest) HasConflicts() bool {
	return p.MergeableState == "dirty" || (p.Mergeable != nil && !*p.Mergeable)
}

type pullPayload struct {
	Number         int    `json:"number"`
	NodeID         string `json:"node_id"`
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
		Ref  string `json:"ref"`
		SHA  string `json:"sha"`
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
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
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &payload); err != nil {
		return PullRequest{}, err
	}

	return payload.toPullRequest(), nil
}

// toPullRequest converts a decoded pull-request payload to the exported type.
func (p pullPayload) toPullRequest() PullRequest {
	return PullRequest{
		Number:         p.Number,
		NodeID:         p.NodeID,
		Title:          p.Title,
		Author:         p.User.Login,
		State:          p.State,
		Draft:          p.Draft,
		HeadRef:        p.Head.Ref,
		HeadSHA:        p.Head.SHA,
		BaseRef:        p.Base.Ref,
		Mergeable:      p.Mergeable,
		MergeableState: p.MergeableState,
		Body:           p.Body,
		URL:            p.HTMLURL,
		HeadRepo:       p.Head.Repo.FullName,
	}
}
