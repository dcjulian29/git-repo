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
)

// DuplicateIssue marks an issue as a duplicate of another: it posts a
// "Duplicate of #<of>" comment (with an optional extra note) and closes the
// issue as "not planned", after confirmation unless skipConfirm is set.
func DuplicateIssue(ctx context.Context, ref Ref, of int, note string, skipConfirm bool) error {
	if of <= 0 {
		return errors.New("--of <number> is required to mark an issue as a duplicate")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	comment := fmt.Sprintf("Duplicate of #%d", of)
	if note != "" {
		comment = comment + "\n\n" + note
	}

	description := fmt.Sprintf("a duplicate of #%d", of)

	return performClose(ctx, client, ref, comment, "not_planned", description, skipConfirm)
}
