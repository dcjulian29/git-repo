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

// Label describes a repository label.
type Label struct {
	// Name is the label name.
	Name string

	// Color is the six-character hex colour, without a leading '#'.
	Color string

	// Description is the optional label description.
	Description string
}

// ListLabels returns every label defined in the repository, following pagination.
func (c *Client) ListLabels(ctx context.Context, repo Repo) ([]Label, error) {
	var labels []Label

	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/%s/labels?per_page=%d&page=%d",
			c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), pageSize, page)

		var batch []struct {
			Name        string `json:"name"`
			Color       string `json:"color"`
			Description string `json:"description"`
		}

		if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &batch); err != nil {
			return nil, err
		}

		for _, label := range batch {
			labels = append(labels, Label{Name: label.Name, Color: label.Color, Description: label.Description})
		}

		if len(batch) < pageSize {
			break
		}
	}

	return labels, nil
}

// Labels returns the names of every label defined in the repository.
func (c *Client) Labels(ctx context.Context, repo Repo) ([]string, error) {
	labels, err := c.ListLabels(ctx, repo)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(labels))
	for i, label := range labels {
		names[i] = label.Name
	}

	return names, nil
}

// CreateLabel creates a new label in the repository.
func (c *Client) CreateLabel(ctx context.Context, repo Repo, label Label) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/labels",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name))

	body := map[string]string{
		"name":        label.Name,
		"color":       label.Color,
		"description": label.Description,
	}

	return c.doJSON(ctx, http.MethodPost, endpoint, body, nil)
}

// DeleteLabel deletes a label from the repository by name.
func (c *Client) DeleteLabel(ctx context.Context, repo Repo, name string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/labels/%s",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), url.PathEscape(name))

	return c.doJSON(ctx, http.MethodDelete, endpoint, nil, nil)
}
