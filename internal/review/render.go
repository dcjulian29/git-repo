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
	"time"
	"unicode/utf8"

	"github.com/dcjulian29/git-repo/internal/shared"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/term"
)

// minTitleWidth is the smallest width the title column is shrunk to before the
// table is allowed to exceed a narrow terminal.
const minTitleWidth = 12

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

// RenderTable writes items as a colour-coded table. The first column combines
// the repository name and number into a "<repo>#<number>" handle (the same
// handle other commands accept), labelled by handleColumn (for example "PR").
// Columns size to their content and never wrap. When there are no items it
// prints emptyMessage instead.
func RenderTable(items []NamedItem, handleColumn, emptyMessage string) error {
	if len(items) == 0 {
		fmt.Println(textformat.Info(emptyMessage))

		return nil
	}

	type row struct {
		handle string
		title  string
		author string
		age    string
	}

	// Track the natural (plain-text) width of every column so the title budget
	// can be derived from what the other columns actually need.
	handleWidth := utf8.RuneCountInString(handleColumn)
	authorWidth := utf8.RuneCountInString("AUTHOR")
	ageWidth := utf8.RuneCountInString("AGE")

	rows := make([]row, len(items))

	for i, item := range items {
		handle := fmt.Sprintf("%s#%d", item.Repo, item.Number)

		title := item.Title
		if item.Draft {
			title += " (draft)"
		}

		age := shared.Age(time.Since(item.CreatedAt))
		rows[i] = row{handle: handle, title: title, author: item.Author, age: age}

		handleWidth = max(handleWidth, utf8.RuneCountInString(handle))
		authorWidth = max(authorWidth, utf8.RuneCountInString(item.Author))
		ageWidth = max(ageWidth, utf8.RuneCountInString(age))
	}

	// When writing to a terminal, cap the title column to the remaining width so
	// the whole table fits on one line; otherwise leave titles intact.
	if width, ok := terminalWidth(); ok {
		titleWidth := titleBudget(width, handleWidth, authorWidth, ageWidth)
		for i := range rows {
			rows[i].title = truncate(rows[i].title, titleWidth)
		}
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.Off},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)

	table.Header([]string{handleColumn, "TITLE", "AUTHOR", "AGE"})

	for _, r := range rows {
		if err := table.Append([]string{color.CyanString(r.handle), r.title, r.author, r.age}); err != nil {
			return err
		}
	}

	return table.Render()
}

// titleBudget returns the width to allow for the title column so that the whole
// table fits within termWidth, never shrinking it below minTitleWidth.
func titleBudget(termWidth, handleWidth, authorWidth, ageWidth int) int {
	const chrome = 3*4 + 1 // padding + separators for four columns

	return max(termWidth-chrome-handleWidth-authorWidth-ageWidth, minTitleWidth)
}

// terminalWidth returns the width of the terminal attached to stdout, or false
// when stdout is not a terminal (for example when the output is piped).
func terminalWidth() (int, bool) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 0, false
	}

	return width, true
}

// truncate shortens s to at most maxWidth runes, appending an ellipsis when it
// has to cut the string.
func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}

	if maxWidth == 1 {
		return "…"
	}

	return string([]rune(s)[:maxWidth-1]) + "…"
}

// RenderJSON writes items as indented JSON to stdout.
func RenderJSON(items []NamedItem) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(items)
}

