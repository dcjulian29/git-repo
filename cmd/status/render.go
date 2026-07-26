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
	"github.com/olekukonko/tablewriter/tw"
)

func render(results []git.RepoStatus) {
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

	table.Header([]string{"PATH", "DIRTY", "PUSH", "PULL", "DIVERGED", "UNTRACKED", "NO UPSTREAM"})

	for _, s := range results {
		if actions {
			_ = table.Append([]string{
				git.ColorPath(s),
				git.ActionLabel(s.Dirty, "dirty", false),
				git.ActionLabel(s.PushNeeded, "push needed", true),
				git.ActionLabel(s.PullNeeded, "pull needed", true),
				git.ActionLabel(s.Diverged, "diverged", true),
				git.ActionLabel(s.Untracked, "untracked files", false),
				git.ActionLabel(s.NoUpstream, "no upstream", false),
			})
		} else {
			_ = table.Append([]string{
				git.ColorPath(s),
				git.ColorBool(s.Dirty, false),
				git.ColorBool(s.PushNeeded, true),
				git.ColorBool(s.PullNeeded, true),
				git.ColorBool(s.Diverged, true),
				git.ColorBool(s.Untracked, false),
				git.ColorBool(s.NoUpstream, false),
			})
		}
	}

	fmt.Println()
	_ = table.Render()

	fmt.Printf("\n  %s = clean   %s = dirty / untracked / no upstream   %s = sync needed\n",
		color.GreenString("■"),
		color.YellowString("■"),
		color.RedString("■"),
	)
}
