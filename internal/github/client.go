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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.github.com"
	pageSize       = 100
	requestTimeout = 30 * time.Second
)

// Client is a minimal GitHub REST client scoped to the endpoints git-repo
// needs. The zero value is not usable; construct one with NewClient.
type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

// NewClient returns a Client that authenticates requests with the given token.
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: requestTimeout},
	}
}

// ListOptions filters the items returned by ListItems. The zero value lists
// open items.
type ListOptions struct {
	// State is "open" (the default when empty), "closed", or "all".
	State string

	// Labels, when set, restricts results to items carrying all of the labels.
	Labels []string

	// Assignee, when set, restricts results to items assigned to that login
	// (or the special values "*" for any and "none" for unassigned).
	Assignee string
}

// ListItems returns every issue and pull request for repo matching opts,
// following pagination. Pull requests are included because the issues endpoint
// returns them; callers separate the two via Item.IsPull.
func (c *Client) ListItems(ctx context.Context, repo Repo, opts ListOptions) ([]Item, error) {
	state := opts.State
	if state == "" {
		state = "open"
	}

	var items []Item

	for page := 1; ; page++ {
		query := url.Values{}
		query.Set("state", state)
		query.Set("per_page", strconv.Itoa(pageSize))
		query.Set("page", strconv.Itoa(page))

		if len(opts.Labels) > 0 {
			query.Set("labels", strings.Join(opts.Labels, ","))
		}

		if opts.Assignee != "" {
			query.Set("assignee", opts.Assignee)
		}

		endpoint := fmt.Sprintf("%s/repos/%s/%s/issues?%s",
			c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), query.Encode())

		var batch []issuePayload
		if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &batch); err != nil {
			return nil, err
		}

		for _, payload := range batch {
			items = append(items, payload.toItem())
		}

		if len(batch) < pageSize {
			break
		}
	}

	return items, nil
}

// doJSON performs an authenticated request against endpoint. When body is
// non-nil it is JSON-encoded as the request payload; when out is non-nil the
// response body is decoded into it. Any non-2xx status is returned as an error.
func (c *Client) doJSON(ctx context.Context, method, endpoint string, body, out any) error {
	var reader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader) //nolint:gosec // G107: endpoint is built from the fixed api base and URL-escaped path
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 512))

		return fmt.Errorf("github api request failed (%s): %s", resp.Status, string(message))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding github response: %w", err)
		}
	}

	return nil
}
