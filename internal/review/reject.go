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

import "context"

// RejectIssue posts a required explanatory comment and then closes the issue as
// "not planned", after confirmation unless skipConfirm is set.
func RejectIssue(ctx context.Context, ref Ref, comment string, skipConfirm bool) error {
	comment, err := requireComment(comment, "Reason for rejecting (required): ")
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	return closeIssue(ctx, client, ref, comment, "not_planned", "not planned", skipConfirm)
}
