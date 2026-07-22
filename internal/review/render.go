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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// PrintWarnings writes skipped (non-github.com) repositories and per-repository
// fetch failures to w, so read commands can surface problems without aborting.
func PrintWarnings(w io.Writer, report Report) {
	for _, name := range report.Skipped {
		_, _ = fmt.Fprintln(w, textformat.Warn(fmt.Sprintf("Skipping %s: not a github.com repository.", name)))
	}

	for _, res := range report.Failures() {
		_, _ = fmt.Fprintln(w, textformat.Warn(fmt.Sprintf("Could not fetch %s: %v", res.Target.Name, res.Err)))
	}
}

// RenderTable writes items as a colour-coded table whose number column is
// labelled numberColumn (for example "PR"). When there are no items it prints
// emptyMessage instead.
func RenderTable(items []NamedItem, numberColumn, emptyMessage string) error {
	if len(items) == 0 {
		fmt.Println(textformat.Info(emptyMessage))

		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"REPO", numberColumn, "TITLE", "AUTHOR", "AGE"})

	for _, item := range items {
		if err := table.Append([]string{
			color.CyanString(item.Repo),
			color.GreenString("#" + strconv.Itoa(item.Number)),
			itemTitle(item),
			item.Author,
			shared.Age(time.Since(item.CreatedAt)),
		}); err != nil {
			return err
		}
	}

	return table.Render()
}

// RenderJSON writes items as indented JSON to stdout.
func RenderJSON(items []NamedItem) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(items)
}

// itemTitle returns the item title, annotating draft pull requests.
func itemTitle(item NamedItem) string {
	if item.Draft {
		return item.Title + " " + color.YellowString("(draft)")
	}

	return item.Title
}
