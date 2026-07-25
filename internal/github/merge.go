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
)

// MergeMethods lists the merge methods GitHub accepts.
var MergeMethods = []string{"merge", "squash", "rebase"}

// MergePull merges the given pull request using method (one of MergeMethods).
func (c *Client) MergePull(ctx context.Context, repo Repo, number int, method string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/merge",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	body, err := json.Marshal(map[string]string{"merge_method": method})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(body)) //nolint:gosec // G107: endpoint is built from the fixed api base and URL-escaped path
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 512))

		return fmt.Errorf("merge failed (%s): %s", resp.Status, string(message))
	}

	return nil
}
