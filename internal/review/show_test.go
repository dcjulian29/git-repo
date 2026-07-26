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
	"strings"
	"testing"

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/fatih/color"
)

func TestAggregateChecks(t *testing.T) {
	pass := func(name string) github.CheckRun {
		return github.CheckRun{Name: name, Status: "completed", Conclusion: "success"}
	}
	fail := func(name string) github.CheckRun {
		return github.CheckRun{Name: name, Status: "completed", Conclusion: "failure"}
	}
	running := func(name string) github.CheckRun {
		return github.CheckRun{Name: name, Status: "in_progress"}
	}

	// lint appears three times (pass, fail, pass); build twice (both pass);
	// test once and still running. First-seen order is lint, build, test.
	checks := []github.CheckRun{
		pass("lint"), pass("build"), fail("lint"), running("test"), pass("build"), pass("lint"),
	}

	got := aggregateChecks(checks)

	if len(got) != 3 {
		t.Fatalf("expected 3 summaries, got %d", len(got))
	}

	want := []struct {
		name       string
		count      int
		conclusion string
		status     string
	}{
		{"lint", 3, "failure", "completed"}, // worst-case run wins
		{"build", 2, "success", "completed"},
		{"test", 1, "", "in_progress"},
	}

	for i, w := range want {
		if got[i].run.Name != w.name {
			t.Fatalf("summary %d name = %q, want %q", i, got[i].run.Name, w.name)
		}

		if got[i].count != w.count {
			t.Fatalf("%s count = %d, want %d", w.name, got[i].count, w.count)
		}

		if got[i].run.Conclusion != w.conclusion || got[i].run.Status != w.status {
			t.Fatalf("%s run = %+v, want status %q conclusion %q",
				w.name, got[i].run, w.status, w.conclusion)
		}
	}
}

func TestCheckSeverity(t *testing.T) {
	fail := github.CheckRun{Name: "x", Status: "completed", Conclusion: "failure"}
	running := github.CheckRun{Name: "x", Status: "in_progress"}
	pass := github.CheckRun{Name: "x", Status: "completed", Conclusion: "success"}

	if checkSeverity(fail) <= checkSeverity(running) || checkSeverity(running) <= checkSeverity(pass) {
		t.Fatalf("severity order wrong: fail=%d running=%d pass=%d",
			checkSeverity(fail), checkSeverity(running), checkSeverity(pass))
	}
}

func TestCheckStatusPassedIsGreen(t *testing.T) {
	previous := color.NoColor
	color.NoColor = false

	t.Cleanup(func() { color.NoColor = previous })

	symbol, _ := checkStatus(github.CheckRun{Status: "completed", Conclusion: "success"})

	if !strings.Contains(symbol, "\x1b[32m") {
		t.Fatalf("passing symbol %q is not green", symbol)
	}

	// The heavy check mark carries a U+FE0E text-presentation selector so the
	// green applies instead of a color emoji.
	if !strings.ContainsRune(symbol, '︎') {
		t.Fatalf("passing symbol %q lacks the text-presentation selector", symbol)
	}
}
