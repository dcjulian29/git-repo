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

package review

import (
	"reflect"
	"testing"

	"github.com/dcjulian29/git-repo/internal/github"
)

func TestCheckoutPlanSameRepo(t *testing.T) {
	ref := Ref{Number: 8, Repo: github.Repo{Owner: "dcjulian29", Name: "app"}}
	pull := github.PullRequest{HeadRepo: "dcjulian29/app", HeadRef: "dependabot/foo"}

	got := checkoutPlan(ref, pull)
	want := [][]string{
		{"fetch", "origin"},
		{"checkout", "dependabot/foo"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("same-repo plan = %v, want %v", got, want)
	}
}

func TestCheckoutPlanFork(t *testing.T) {
	ref := Ref{Number: 8, Repo: github.Repo{Owner: "dcjulian29", Name: "app"}}
	pull := github.PullRequest{HeadRepo: "someone-else/app", HeadRef: "feature"}

	got := checkoutPlan(ref, pull)
	want := [][]string{
		{"fetch", "origin", "pull/8/head:pr/8"},
		{"checkout", "pr/8"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fork plan = %v, want %v", got, want)
	}
}
