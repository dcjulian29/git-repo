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

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type (
	// CheckRun is a single CI check on a commit.
	CheckRun struct {
		// Name is the check's display name.
		Name string

		// Status is the run status: "queued", "in_progress", or "completed".
		Status string

		// Conclusion is the result once completed: "success", "failure",
		// "neutral", "cancelled", "timed_out", "action_required", or "skipped".
		// It is empty while the check is still running.
		Conclusion string
	}

	// checkRunPayload mirrors a single entry in the check-runs API response.
	checkRunPayload struct {
		Name       string `json:"name"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}
)

// Completed reports whether the check has finished running.
func (r CheckRun) Completed() bool {
	return r.Status == "completed"
}

// Passed reports whether a completed check succeeded (or was skipped/neutral).
func (r CheckRun) Passed() bool {
	switch r.Conclusion {
	case "success", "neutral", "skipped":
		return true
	default:
		return false
	}
}

// Checks returns the check runs for the given commit SHA.
func (c *Client) Checks(ctx context.Context, repo Repo, sha string) ([]CheckRun, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs?per_page=%d",
		c.baseURL, url.PathEscape(repo.Owner), url.PathEscape(repo.Name), url.PathEscape(sha), pageSize)

	var payload struct {
		CheckRuns []checkRunPayload `json:"check_runs"`
	}

	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}

	runs := make([]CheckRun, len(payload.CheckRuns))
	for i, run := range payload.CheckRuns {
		runs[i] = CheckRun(run)
	}

	return runs, nil
}
