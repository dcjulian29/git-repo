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

// Issue holds the details git-repo needs to triage a single issue.
type Issue struct {
	// Number is the issue number.
	Number int

	// Title is the issue title.
	Title string

	// Author is the login of the user who opened it.
	Author string

	// State is the issue state (for example "open").
	State string

	// Body is the issue description.
	Body string

	// Labels are the names of the labels currently applied.
	Labels []string

	// Assignees are the logins currently assigned.
	Assignees []string

	// Comments is the number of comments on the issue.
	Comments int

	// CreatedAt is the time the issue was opened.
	CreatedAt time.Time

	// URL is the html_url that opens the issue in a browser.
	URL string
}

type issueDetailPayload struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Body      string    `json:"body"`
	Comments  int       `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	HTMLURL   string    `json:"html_url"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignees []struct {
		Login string `json:"login"`
	} `json:"assignees"`
}

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

// Labels returns the names of every label defined in the repository, following
// pagination.
func (c *Client) Labels(ctx context.Context, repo Repo) ([]string, error) {
	var names []string

	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/%s/labels?per_page=%d&page=%d",
			c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), pageSize, page)

		var batch []struct {
			Name string `json:"name"`
		}

		if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &batch); err != nil {
			return nil, err
		}

		for _, label := range batch {
			names = append(names, label.Name)
		}

		if len(batch) < pageSize {
			break
		}
	}

	return names, nil
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
