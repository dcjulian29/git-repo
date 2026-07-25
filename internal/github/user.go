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
	"net/http"
)

// Viewer returns the login of the authenticated user.
func (c *Client) Viewer(ctx context.Context) (string, error) {
	var payload struct {
		Login string `json:"login"`
	}

	if err := c.doJSON(ctx, http.MethodGet, c.baseURL+"/user", nil, &payload); err != nil {
		return "", err
	}

	return payload.Login, nil
}
