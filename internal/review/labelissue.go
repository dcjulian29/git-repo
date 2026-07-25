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
	"strings"

	"github.com/fatih/color"
)

// LabelIssue adds one or more labels to an issue. Each label is validated
// against the repository's defined labels before any is applied.
func LabelIssue(ctx context.Context, ref Ref, labels []string) error {
	if len(labels) == 0 {
		return errors.New("at least one label is required")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	available, err := client.Labels(ctx, ref.Repo)
	if err != nil {
		return err
	}

	canonical := make([]string, 0, len(labels))

	for _, want := range labels {
		match, ok := matchLabel(available, want)
		if !ok {
			return fmt.Errorf("label %q does not exist in %s (available: %s)",
				want, ref.Name, strings.Join(available, ", "))
		}

		canonical = append(canonical, match)
	}

	if err := client.AddLabels(ctx, ref.Repo, ref.Number, canonical); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Added %s to %s#%d.",
		strings.Join(canonical, ", "), ref.Name, ref.Number))

	return nil
}
