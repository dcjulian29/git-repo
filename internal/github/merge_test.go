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

func TestMergePullSuccess(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}

		if r.URL.Path != "/repos/o/app/pulls/8/merge" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)

		var payload map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("bad request body: %v", err)
		}

		if payload["merge_method"] != "squash" {
			t.Errorf("merge_method = %q, want squash", payload["merge_method"])
		}

		_, _ = w.Write([]byte(`{"merged":true,"message":"Pull Request successfully merged"}`))
	})

	if err := client.MergePull(context.Background(), Repo{Owner: "o", Name: "app"}, 8, "squash"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMergePullReportsFailure(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"message":"Pull Request is not mergeable"}`))
	})

	err := client.MergePull(context.Background(), Repo{Owner: "o", Name: "app"}, 8, "merge")
	if err == nil {
		t.Fatal("expected an error when the merge is rejected")
	}
}
