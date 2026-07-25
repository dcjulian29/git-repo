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
	"testing"
)

func TestGetIssueParsesLabelsAndAssignees(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/o/app/issues/5" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{
			"number":5,"title":"broken","state":"open","comments":3,
			"user":{"login":"alice"},
			"labels":[{"name":"bug"},{"name":"help wanted"}],
			"assignees":[{"login":"bob"}]
		}`))
	})

	issue, err := client.GetIssue(context.Background(), Repo{Owner: "o", Name: "app"}, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if issue.Number != 5 || issue.Author != "alice" || issue.Comments != 3 {
		t.Fatalf("issue metadata wrong: %+v", issue)
	}

	if len(issue.Labels) != 2 || issue.Labels[0] != "bug" || issue.Labels[1] != "help wanted" {
		t.Fatalf("labels = %v", issue.Labels)
	}

	if len(issue.Assignees) != 1 || issue.Assignees[0] != "bob" {
		t.Fatalf("assignees = %v", issue.Assignees)
	}
}

func TestLabelsListsNames(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/o/app/labels" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(`[{"name":"bug"},{"name":"accepted"}]`))
	})

	names, err := client.Labels(context.Background(), Repo{Owner: "o", Name: "app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(names) != 2 || names[1] != "accepted" {
		t.Fatalf("labels = %v", names)
	}
}

func TestAddLabelsSendsPost(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/o/app/issues/5/labels" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)

		var payload map[string][]string
		_ = json.Unmarshal(body, &payload)

		if len(payload["labels"]) != 1 || payload["labels"][0] != "accepted" {
			t.Errorf("labels body = %v", payload)
		}

		_, _ = w.Write([]byte(`[]`))
	})

	if err := client.AddLabels(context.Background(), Repo{Owner: "o", Name: "app"}, 5, []string{"accepted"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloseIssueSendsPatch(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}

		body, _ := io.ReadAll(r.Body)

		var payload map[string]string
		_ = json.Unmarshal(body, &payload)

		if payload["state"] != "closed" || payload["state_reason"] != "completed" {
			t.Errorf("close body = %v", payload)
		}

		_, _ = w.Write([]byte(`{"number":5,"state":"closed"}`))
	})

	if err := client.CloseIssue(context.Background(), Repo{Owner: "o", Name: "app"}, 5, "completed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
