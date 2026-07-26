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
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMarkPullReadySendsMutation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/graphql" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)

		var payload struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		_ = json.Unmarshal(body, &payload)

		if !strings.Contains(payload.Query, "markPullRequestReadyForReview") {
			t.Errorf("query missing mutation: %s", payload.Query)
		}

		if payload.Variables["id"] != "PR_node_123" {
			t.Errorf("variables id = %v, want PR_node_123", payload.Variables["id"])
		}

		_, _ = w.Write([]byte(`{"data":{"markPullRequestReadyForReview":{"pullRequest":{"number":8}}}}`))
	})

	if err := client.MarkPullReady(context.Background(), "PR_node_123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarkPullReadySurfacesGraphQLErrors(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		// GraphQL reports application errors with a 200 status.
		_, _ = w.Write([]byte(`{"errors":[{"message":"Pull request is not a draft"}]}`))
	})

	err := client.MarkPullReady(context.Background(), "PR_node_123")
	if err == nil {
		t.Fatal("expected an error from a GraphQL error response")
	}

	if !strings.Contains(err.Error(), "not a draft") {
		t.Fatalf("error should include the GraphQL message: %v", err)
	}
}
