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
	"errors"
	"fmt"
	"os"

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/fatih/color"
)

// closeIssue shows the issue, confirms (unless skipConfirm), posts comment when
// it is non-empty, and closes the issue with the given state reason. description
// appears in the confirmation prompt (for example "completed").
func closeIssue(
	ctx context.Context,
	client *github.Client,
	ref Ref,
	comment, reason, description string,
	skipConfirm bool,
) error {
	issue, err := client.GetIssue(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	handle := color.CyanString("%s#%d", ref.Name, ref.Number)
	fmt.Printf("%s %s\n", handle, issue.Title)

	if !skipConfirm {
		confirmed, err := shared.Confirm(os.Stdin, os.Stdout,
			fmt.Sprintf("Close %s as %s? [y/N] ", handle, description))
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Aborted.")

			return nil
		}
	}

	if comment != "" {
		if err := client.CommentIssue(ctx, ref.Repo, ref.Number, comment); err != nil {
			return err
		}
	}

	if err := client.CloseIssue(ctx, ref.Repo, ref.Number, reason); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Closed %s#%d.", ref.Name, ref.Number))

	return nil
}

// requireComment returns comment, prompting for it interactively when empty. It
// returns an error when no comment is ultimately provided.
func requireComment(comment, prompt string) (string, error) {
	if comment == "" {
		answer, err := shared.Prompt(os.Stdin, os.Stdout, prompt)
		if err != nil {
			return "", err
		}

		comment = answer
	}

	if comment == "" {
		return "", errors.New("a comment is required")
	}

	return comment, nil
}
