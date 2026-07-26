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

// RenderTable writes items as a color-coded table. The first column combines
// the repository name and number into a "<repo>#<number>" handle (the same
// handle other commands accept), labeled by handleColumn (for example "PR").
// When showMerge is set, a trailing MERGE column shows each pull request's merge
// state, colored green for "clean" and red otherwise. Columns size to their
// content and never wrap; when there are no items it prints emptyMessage.
func RenderTable(items []NamedItem, handleColumn, emptyMessage string, showMerge bool) error {
	if len(items) == 0 {
		fmt.Println(textformat.Info(emptyMessage))

		return nil
	}

	type row struct {
		handle string
		title  string
		author string
		age    string
		merge  string
	}

	// Track the natural (plain-text) width of every column so the title budget
	// can be derived from what the other columns actually need.
	handleWidth := utf8.RuneCountInString(handleColumn)
	authorWidth := utf8.RuneCountInString("AUTHOR")
	ageWidth := utf8.RuneCountInString("AGE")
	mergeWidth := 0

	if showMerge {
		mergeWidth = utf8.RuneCountInString("MERGE")
	}

	rows := make([]row, len(items))
	anyMergeState := false

	for i, item := range items {
		handle := fmt.Sprintf("%s#%d", item.Repo, item.Number)

		title := item.Title
		if item.Draft {
			title += " (draft)"
		}

		age := shared.Age(time.Since(item.CreatedAt))
		rows[i] = row{handle: handle, title: title, author: item.Author, age: age, merge: item.MergeableState}

		handleWidth = max(handleWidth, utf8.RuneCountInString(handle))
		authorWidth = max(authorWidth, utf8.RuneCountInString(item.Author))
		ageWidth = max(ageWidth, utf8.RuneCountInString(age))

		if item.MergeableState != "" {
			anyMergeState = true
		}

		if showMerge {
			mergeWidth = max(mergeWidth, utf8.RuneCountInString(mergeStatePlain(item.MergeableState)))
		}
	}

	// When writing to a terminal, cap the title column to the remaining width so
	// the whole table fits on one line; otherwise leave titles intact.
	if width, ok := terminalWidth(); ok {
		columns := 4
		otherWidth := handleWidth + authorWidth + ageWidth

		if showMerge {
			columns = 5
			otherWidth += mergeWidth
		}

		titleWidth := titleBudget(width, otherWidth, columns)
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

	header := []string{handleColumn, "TITLE", "AUTHOR", "AGE"}
	if showMerge {
		header = append(header, "MERGE")
	}

	table.Header(header)

	for _, r := range rows {
		cells := []string{colorHandle(r.handle, r.merge), r.title, r.author, r.age}
		if showMerge {
			cells = append(cells, colorMergeState(r.merge))
		}

		if err := table.Append(cells); err != nil {
			return err
		}
	}

	if err := table.Render(); err != nil {
		return err
	}

	if showMerge && anyMergeState {
		fmt.Printf("\n  %s = clean   %s = unknown   %s = unstable\n",
			color.GreenString("■"), color.YellowString("■"), color.RedString("■"))
	}

	return nil
}

// titleBudget returns the width to allow for the title column so that the whole
// table fits within termWidth, never shrinking it below minTitleWidth.
// otherWidth is the combined content width of the non-title columns and columns
// is the total column count.
func titleBudget(termWidth, otherWidth, columns int) int {
	chrome := 3*columns + 1 // padding + separators

	return max(termWidth-chrome-otherWidth, minTitleWidth)
}

// mergeStatePlain returns the plain display text for a merge state, using "-"
// for an unknown or empty state.
func mergeStatePlain(state string) string {
	if state == "" {
		return "-"
	}

	return state
}

// colorHandle colors a "<repo>#<number>" handle by merge state: green when
// clean, red for problem states such as "unstable" or "dirty", yellow when
// unknown, and the neutral cyan when there is no merge state (issues, or pull
// requests whose merge state was not fetched).
func colorHandle(handle, state string) string {
	switch state {
	case "":
		return color.CyanString(handle)
	case "clean":
		return color.GreenString(handle)
	case "unstable", "dirty", "blocked", "behind":
		return color.RedString(handle)
	default:
		return color.YellowString(handle)
	}
}

// colorMergeState returns the merge state colored green when clean, yellow when
// unknown, and red for any problem state such as "unstable" or "dirty".
func colorMergeState(state string) string {
	switch state {
	case "":
		return "-"
	case "clean":
		return color.GreenString("clean")
	case "unknown":
		return color.YellowString("unknown")
	default:
		return color.RedString(state)
	}
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

