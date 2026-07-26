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
	"encoding/json"
	"fmt"
	"net/http"
)

// graphQL executes a GraphQL query against the GitHub API. GraphQL reports
// application errors in the response body with a 200 status, so those are
// surfaced as errors here. When out is non-nil the "data" object is decoded
// into it.
func (c *Client) graphQL(ctx context.Context, query string, variables map[string]any, out any) error {
	request := map[string]any{"query": query, "variables": variables}

	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.doJSON(ctx, http.MethodPost, c.baseURL+"/graphql", request, &envelope); err != nil {
		return err
	}

	if len(envelope.Errors) > 0 {
		return fmt.Errorf("github graphql: %s", envelope.Errors[0].Message)
	}

	if out != nil {
		return json.Unmarshal(envelope.Data, out)
	}

	return nil
}

// markReadyMutation marks a draft pull request ready for review. Toggling the
// draft state is only available through GraphQL, not the REST API.
const markReadyMutation = `mutation($id: ID!) {
  markPullRequestReadyForReview(input: {pullRequestId: $id}) {
    pullRequest { number }
  }
}`

// MarkPullReady marks the pull request with the given node ID ready for review.
func (c *Client) MarkPullReady(ctx context.Context, nodeID string) error {
	return c.graphQL(ctx, markReadyMutation, map[string]any{"id": nodeID}, nil)
}
