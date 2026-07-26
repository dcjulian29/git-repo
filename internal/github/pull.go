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

type (
	// PullRequest holds the details git-repo needs to review a single pull request.
	PullRequest struct {
		Mergeable      *bool
		Author         string
		Title          string
		State          string
		HeadRef        string
		HeadSHA        string
		BaseRef        string
		NodeID         string
		MergeableState string
		Body           string
		URL            string
		HeadRepo       string
		Number         int
		Draft          bool
	}

	// headPayload is the head branch reference of a pull request.
	headPayload struct {
		Ref  string      `json:"ref"`
		SHA  string      `json:"sha"`
		Repo repoPayload `json:"repo"`
	}

	// repoPayload identifies the repository a branch lives in.
	repoPayload struct {
		FullName string `json:"full_name"`
	}

	// basePayload is the base branch reference of a pull request.
	basePayload struct {
		Ref string `json:"ref"`
	}

	pullPayload struct {
		Mergeable      *bool       `json:"mergeable"`
		Head           headPayload `json:"head"`
		NodeID         string      `json:"node_id"`
		Title          string      `json:"title"`
		State          string      `json:"state"`
		Body           string      `json:"body"`
		HTMLURL        string      `json:"html_url"`
		MergeableState string      `json:"mergeable_state"`
		User           userPayload `json:"user"`
		Base           basePayload `json:"base"`
		Number         int         `json:"number"`
		Draft          bool        `json:"draft"`
	}
)

// HasConflicts reports whether the pull request has merge conflicts.
func (p PullRequest) HasConflicts() bool {
	return p.MergeableState == "dirty" || (p.Mergeable != nil && !*p.Mergeable)
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
