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

func TestListLabels(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/o/app/labels" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(`[
			{"name":"bug","color":"d73a4a","description":"Something is broken"},
			{"name":"good first issue","color":"7057ff","description":""}
		]`))
	})

	labels, err := client.ListLabels(context.Background(), Repo{Owner: "o", Name: "app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(labels) != 2 {
		t.Fatalf("got %d labels, want 2", len(labels))
	}

	if labels[0].Name != "bug" || labels[0].Color != "d73a4a" || labels[0].Description != "Something is broken" {
		t.Fatalf("first label = %+v", labels[0])
	}
}

func TestCreateLabelSendsPost(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/o/app/labels" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)

		var payload map[string]string
		_ = json.Unmarshal(body, &payload)

		if payload["name"] != "accepted" || payload["color"] != "0e8a16" {
			t.Errorf("create body = %v", payload)
		}

		_, _ = w.Write([]byte(`{"name":"accepted"}`))
	})

	err := client.CreateLabel(context.Background(), Repo{Owner: "o", Name: "app"},
		Label{Name: "accepted", Color: "0e8a16", Description: "Accepted for work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteLabelEscapesName(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}

		// A label name with a space must be percent-encoded in the path.
		if r.URL.EscapedPath() != "/repos/o/app/labels/good%20first%20issue" {
			t.Errorf("unexpected escaped path %q", r.URL.EscapedPath())
		}

		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteLabel(context.Background(), Repo{Owner: "o", Name: "app"}, "good first issue"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
