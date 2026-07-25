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

func TestDefaultBranch(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/o/app" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"default_branch":"main"}`))
	})

	branch, err := client.DefaultBranch(context.Background(), Repo{Owner: "o", Name: "app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if branch != "main" {
		t.Fatalf("DefaultBranch = %q, want main", branch)
	}
}

func TestCreatePull(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/o/app/pulls" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)

		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if payload["head"] != "issue/8" || payload["base"] != "main" || payload["draft"] != true {
			t.Errorf("create body = %v", payload)
		}

		if payload["body"] != "Closes #8" {
			t.Errorf("body should link the issue, got %v", payload["body"])
		}

		_, _ = w.Write([]byte(`{"number":21,"html_url":"https://gh/pull/21"}`))
	})

	pull, err := client.CreatePull(context.Background(), Repo{Owner: "o", Name: "app"}, CreatePullParams{
		Title: "Fix it",
		Head:  "issue/8",
		Base:  "main",
		Body:  "Closes #8",
		Draft: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pull.Number != 21 || pull.URL != "https://gh/pull/21" {
		t.Fatalf("created pull = %+v", pull)
	}
}
