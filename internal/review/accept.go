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
	"strings"

	"github.com/fatih/color"
)

// AcceptIssue marks an issue accepted by applying an existing label (validated
// against the repository's labels) and optionally assigning it. A label is
// required; there is no default, since repositories define no "accepted" label
// by convention.
func AcceptIssue(ctx context.Context, ref Ref, label, assignee string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	available, err := client.Labels(ctx, ref.Repo)
	if err != nil {
		return err
	}

	if label == "" {
		return fmt.Errorf("a label is required; pass --label (available: %s)", strings.Join(available, ", "))
	}

	canonical, ok := matchLabel(available, label)
	if !ok {
		return fmt.Errorf("label %q does not exist in %s (available: %s)",
			label, ref.Name, strings.Join(available, ", "))
	}

	if err := client.AddLabels(ctx, ref.Repo, ref.Number, []string{canonical}); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Labeled %s#%d as %q.", ref.Name, ref.Number, canonical))

	if assignee != "" {
		if err := client.AddAssignees(ctx, ref.Repo, ref.Number, []string{assignee}); err != nil {
			return err
		}

		fmt.Println(color.GreenString("Assigned to %s.", assignee))
	}

	return nil
}
