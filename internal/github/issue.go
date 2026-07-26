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
	"time"
)

type (
	// Issue holds the details git-repo needs to triage a single issue.
	Issue struct {
		CreatedAt time.Time
		Title     string
		Author    string
		State     string
		Body      string
		URL       string
		Labels    []string
		Assignees []string
		Number    int
		Comments  int
	}

	// userPayload is a GitHub user reference carrying just the login.
	userPayload struct {
		Login string `json:"login"`
	}

	// labelPayload is a label reference carrying just the name.
	labelPayload struct {
		Name string `json:"name"`
	}

	issueDetailPayload struct {
		CreatedAt time.Time      `json:"created_at"`
		Title     string         `json:"title"`
		State     string         `json:"state"`
		Body      string         `json:"body"`
		HTMLURL   string         `json:"html_url"`
		User      userPayload    `json:"user"`
		Labels    []labelPayload `json:"labels"`
		Assignees []userPayload  `json:"assignees"`
		Number    int            `json:"number"`
		Comments  int            `json:"comments"`
	}
)

// GetIssue fetches the details of a single issue.
func (c *Client) GetIssue(ctx context.Context, repo Repo, number int) (Issue, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	var payload issueDetailPayload
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &payload); err != nil {
		return Issue{}, err
	}

	issue := Issue{
		Number:    payload.Number,
		Title:     payload.Title,
		Author:    payload.User.Login,
		State:     payload.State,
		Body:      payload.Body,
		Comments:  payload.Comments,
		CreatedAt: payload.CreatedAt,
		URL:       payload.HTMLURL,
	}

	for _, label := range payload.Labels {
		issue.Labels = append(issue.Labels, label.Name)
	}

	for _, assignee := range payload.Assignees {
		issue.Assignees = append(issue.Assignees, assignee.Login)
	}

	return issue, nil
}

// AddLabels adds the given labels to an issue.
func (c *Client) AddLabels(ctx context.Context, repo Repo, number int, labels []string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d/labels",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	return c.doJSON(ctx, http.MethodPost, endpoint, map[string][]string{"labels": labels}, nil)
}

// AddAssignees assigns the given logins to an issue.
func (c *Client) AddAssignees(ctx context.Context, repo Repo, number int, assignees []string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d/assignees",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	return c.doJSON(ctx, http.MethodPost, endpoint, map[string][]string{"assignees": assignees}, nil)
}

// CommentIssue posts a comment on an issue.
func (c *Client) CommentIssue(ctx context.Context, repo Repo, number int, body string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	return c.doJSON(ctx, http.MethodPost, endpoint, map[string]string{"body": body}, nil)
}

// CloseIssue closes an issue with the given state reason, which is "completed"
// or "not_planned".
func (c *Client) CloseIssue(ctx context.Context, repo Repo, number int, reason string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues/%d",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	payload := map[string]string{"state": "closed", "state_reason": reason}

	return c.doJSON(ctx, http.MethodPatch, endpoint, payload, nil)
}
