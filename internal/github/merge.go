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

// MergeMethods lists the merge methods GitHub accepts.
//
//nolint:gochecknoglobals // exported list of allowed merge methods; a slice cannot be const
var MergeMethods = []string{"merge", "squash", "rebase"}

// MergePull merges the given pull request using method (one of MergeMethods).
func (c *Client) MergePull(ctx context.Context, repo Repo, number int, method string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/merge",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), number)

	return c.doJSON(ctx, http.MethodPut, endpoint, map[string]string{"merge_method": method}, nil)
}
