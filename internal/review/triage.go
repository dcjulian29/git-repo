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
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/fatih/color"
)

// TriageIssue shows an issue and applies a label and/or assignee. When the label
// or assignee is not supplied it is prompted for interactively. A label, however
// supplied, is validated against the repository's defined labels.
func TriageIssue(ctx context.Context, ref Ref, label, assignee string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	renderIssue(ref, issue)

	available, err := client.Labels(ctx, ref.Repo)
	if err != nil {
		return err
	}

	label, err = resolveLabel(label, available)
	if err != nil {
		return err
	}

	if label != "" {
		if err := client.AddLabels(ctx, ref.Repo, ref.Number, []string{label}); err != nil {
			return err
		}

		fmt.Println(color.GreenString("Added label %q.", label))
	}

	if assignee == "" {
		assignee, err = shared.Prompt(os.Stdin, os.Stdout, "Assignee (login, blank to skip): ")
		if err != nil {
			return err
		}
	}

	if assignee != "" {
		if err := client.AddAssignees(ctx, ref.Repo, ref.Number, []string{assignee}); err != nil {
			return err
		}

		fmt.Println(color.GreenString("Assigned to %s.", assignee))
	}

	return nil
}

// resolveLabel returns a validated canonical label. When want is empty it prompts
// interactively (listing the available labels); a blank answer skips labelling.
// A non-empty label must exist among available.
func resolveLabel(want string, available []string) (string, error) {
	if want == "" {
		fmt.Printf("Available labels: %s\n", strings.Join(available, ", "))

		answer, err := shared.Prompt(os.Stdin, os.Stdout, "Label (blank to skip): ")
		if err != nil {
			return "", err
		}

		want = answer
	}

	if want == "" {
		return "", nil
	}

	canonical, ok := matchLabel(available, want)
	if !ok {
		return "", fmt.Errorf("label %q does not exist; available: %s", want, strings.Join(available, ", "))
	}

	return canonical, nil
}

// matchLabel returns the canonically-cased label matching want, if present.
func matchLabel(available []string, want string) (string, bool) {
	for _, name := range available {
		if strings.EqualFold(name, want) {
			return name, true
		}
	}

	return "", false
}

// renderIssue prints an issue's metadata and a body excerpt.
func renderIssue(ref Ref, issue github.Issue) {
	fmt.Printf("%s %s\n", color.CyanString("%s#%d", ref.Name, issue.Number), issue.Title)
	fmt.Printf("  %-11s %s\n", "Author:", issue.Author)
	fmt.Printf("  %-11s %s\n", "State:", issue.State)
	fmt.Printf("  %-11s %s\n", "Age:", shared.Age(time.Since(issue.CreatedAt)))
	fmt.Printf("  %-11s %s\n", "Labels:", joinOrNone(issue.Labels))
	fmt.Printf("  %-11s %s\n", "Assignees:", joinOrNone(issue.Assignees))
	fmt.Printf("  %-11s %d\n", "Comments:", issue.Comments)

	if body := strings.TrimSpace(issue.Body); body != "" {
		fmt.Printf("\n%s\n", excerpt(body))
	}
}

// joinOrNone joins values with commas, or returns "none" when empty.
func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "none"
	}

	return strings.Join(values, ", ")
}

// excerpt truncates long bodies to keep the triage view compact.
func excerpt(text string) string {
	const limit = 600

	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}

	return string(runes[:limit]) + "…"
}
