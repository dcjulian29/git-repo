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

	"github.com/dcjulian29/git-repo/internal/github"
	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// defaultLabelColor is used when a new label is created without a color.
const defaultLabelColor = "ededed"

// ListLabels prints the labels defined in the named repository as a table.
func ListLabels(ctx context.Context, repoName string) error {
	name, repo, err := ResolveRepo(repoName)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	labels, err := client.ListLabels(ctx, repo)
	if err != nil {
		return err
	}

	if len(labels) == 0 {
		fmt.Println(textformat.Info(fmt.Sprintf("No labels defined in %s.", name)))

		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"NAME", "COLOR", "DESCRIPTION"})

	for _, label := range labels {
		row := []string{color.CyanString(label.Name), "#" + label.Color, label.Description}
		if err := table.Append(row); err != nil {
			return err
		}
	}

	return table.Render()
}

// CreateLabel creates a label in the named repository. When color is empty a
// neutral default is used; a leading '#' is accepted and stripped.
func CreateLabel(ctx context.Context, repoName, labelName, labelColor, description string) error {
	_, repo, err := ResolveRepo(repoName)
	if err != nil {
		return err
	}

	if labelColor == "" {
		labelColor = defaultLabelColor
	}

	labelColor = strings.TrimPrefix(labelColor, "#")

	client, err := newClient()
	if err != nil {
		return err
	}

	label := github.Label{Name: labelName, Color: labelColor, Description: description}
	if err := client.CreateLabel(ctx, repo, label); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Created label %q in %s.", labelName, repoName))

	return nil
}

// RemoveLabel deletes a label from the named repository after confirmation
// unless skipConfirm is set.
func RemoveLabel(ctx context.Context, repoName, labelName string, skipConfirm bool) error {
	name, repo, err := ResolveRepo(repoName)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if !skipConfirm {
		confirmed, err := shared.Confirm(os.Stdin, os.Stdout,
			fmt.Sprintf("Delete label %q from %s? [y/N] ", labelName, name))
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Aborted.")

			return nil
		}
	}

	if err := client.DeleteLabel(ctx, repo, labelName); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Deleted label %q from %s.", labelName, name))

	return nil
}
