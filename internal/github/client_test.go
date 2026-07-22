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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := NewClient("test-token")
	client.baseURL = server.URL

	return client
}

func issuesJSON(t *testing.T, count, start int) []byte {
	t.Helper()

	type user struct {
		Login string `json:"login"`
	}

	type item struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		User   user   `json:"user"`
	}

	items := make([]item, count)
	for i := range items {
		items[i] = item{Number: start + i, Title: "generated", User: user{Login: "bob"}}
	}

	data, err := json.Marshal(items)
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}

	return data
}

func TestOpenItemsSeparatesIssuesAndPulls(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want bearer token", got)
		}

		if r.URL.Path != "/repos/dcjulian29/git-repo/issues" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(`[
			{"number":1,"title":"a bug","html_url":"https://gh/1","user":{"login":"alice"}},
			{"number":2,"title":"bump dep","html_url":"https://gh/2","user":{"login":"dependabot[bot]"},"draft":true,"pull_request":{"url":"https://api/2"}}
		]`))
	})

	items, err := client.OpenItems(context.Background(), Repo{Owner: "dcjulian29", Name: "git-repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}

	var issue, pull *Item

	for i := range items {
		if items[i].IsPull {
			pull = &items[i]
		} else {
			issue = &items[i]
		}
	}

	if issue == nil || issue.Number != 1 || issue.Author != "alice" {
		t.Fatalf("issue not parsed correctly: %+v", issue)
	}

	if pull == nil || pull.Number != 2 || pull.Author != "dependabot[bot]" || !pull.Draft {
		t.Fatalf("pull request not parsed correctly: %+v", pull)
	}
}

func TestOpenItemsPaginates(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "1":
			_, _ = w.Write(issuesJSON(t, pageSize, 1))
		case "2":
			_, _ = w.Write(issuesJSON(t, 5, pageSize+1))
		default:
			_, _ = w.Write([]byte(`[]`))
		}
	})

	items, err := client.OpenItems(context.Background(), Repo{Owner: "o", Name: "r"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != pageSize+5 {
		t.Fatalf("got %d items, want %d (pagination not followed)", len(items), pageSize+5)
	}
}

func TestOpenItemsReturnsErrorOnBadStatus(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	})

	_, err := client.OpenItems(context.Background(), Repo{Owner: "o", Name: "missing"})
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("error should mention the status: %v", err)
	}
}

func TestGatherAggregatesAndRecordsErrors(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/good/") {
			_, _ = w.Write([]byte(`[{"number":1,"title":"x","user":{"login":"a"}}]`))

			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	})

	targets := []Target{
		{Name: "good", Repo: Repo{Owner: "o", Name: "good"}},
		{Name: "bad", Repo: Repo{Owner: "o", Name: "bad"}},
	}

	results := client.Gather(context.Background(), targets, 2)
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	if results[0].Err != nil || len(results[0].Items) != 1 {
		t.Fatalf("good target result wrong: %+v", results[0])
	}

	if results[1].Err == nil {
		t.Fatalf("bad target should have recorded an error: %+v", results[1])
	}
}
