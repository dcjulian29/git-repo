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

package status

import (
	"fmt"
	"os"

	"github.com/dcjulian29/git-repo/internal/git"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func render(results []git.RepoStatus) {
	var table *tablewriter.Table

	if actions {
		table = tablewriter.NewTable(os.Stdout,
			tablewriter.WithConfig(tablewriter.Config{
				Row: tw.CellConfig{
					Formatting: tw.CellFormatting{
						AutoWrap:   tw.WrapNone,
						AutoFormat: tw.Off,
					},
					Alignment: tw.CellAlignment{
						Global: tw.AlignLeft,
					},
					Padding: tw.CellPadding{
						Global: tw.Padding{Left: "  ", Right: ""},
					},
				},
				Header: tw.CellConfig{
					Formatting: tw.CellFormatting{
						AutoFormat: tw.Off,
					},
					Padding: tw.CellPadding{
						Global: tw.Padding{Left: "  ", Right: ""},
					},
				},
				Behavior: tw.Behavior{
					TrimSpace: tw.Off,
				},
			}),
			tablewriter.WithRenderer(renderer.NewBlueprint(
				tw.Rendition{

					Borders: tw.Border{
						Left:   tw.Off,
						Top:    tw.Off,
						Right:  tw.Off,
						Bottom: tw.Off,
					},
					Settings: tw.Settings{
						Lines: tw.Lines{
							ShowHeaderLine: tw.Off,
						},
						Separators: tw.Separators{
							BetweenColumns: tw.Off,
							BetweenRows:    tw.Off,
						},
					},
					Symbols: tw.NewSymbols(tw.StyleNone),
				},
			)),
		)
	} else {
		table = tablewriter.NewTable(os.Stdout,
			tablewriter.WithRenderer(renderer.NewBlueprint(
				tw.Rendition{
					Borders: tw.Border{
						Left:   tw.On,
						Top:    tw.Off,
						Right:  tw.On,
						Bottom: tw.Off,
					},
					Settings: tw.Settings{
						Lines: tw.Lines{
							ShowHeaderLine: tw.On,
						},
						Separators: tw.Separators{},
					},
				},
			)),
		)
	}

	for _, s := range results {
		if actions {
			_ = table.Append([]string{
				git.ColorPath(s),
				git.ActionLabel(s.Dirty, "dirty", false),
				git.ActionLabel(s.PushNeeded, "push needed", true),
				git.ActionLabel(s.PullNeeded, "pull needed", true),
				git.ActionLabel(s.Diverged, "diverged", true),
				git.ActionLabel(s.Untracked, "untracked files", false),
			})
		} else {
			_ = table.Append([]string{
				git.ColorPath(s),
				git.ColorBool(s.Dirty, false),
				git.ColorBool(s.PushNeeded, true),
				git.ColorBool(s.PullNeeded, true),
				git.ColorBool(s.Diverged, true),
				git.ColorBool(s.Untracked, false),
			})

			table.Header([]string{"PATH", "DIRTY", "PUSH", "PULL", "DIVERGED", "UNTRACKED"})

		}
	}

	fmt.Println()
	_ = table.Render()

	fmt.Printf("\n  %s = clean   %s = dirty / untracked   %s = sync needed\n",
		color.GreenString("■"),
		color.YellowString("■"),
		color.RedString("■"),
	)
}
