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
	"io"
	"net/http"
	"regexp"
	"time"
)

var (
	compatBadgeExpr = regexp.MustCompile(`https://dependabot-badges\.githubapp\.com/badges/compatibility_score\?[^\s)"']+`)
	compatScoreExpr = regexp.MustCompile(`(\d+)%`)
)

// CompatibilityScore extracts the Dependabot compatibility score referenced in a
// pull-request body, on a best-effort basis. GitHub exposes the score only as a
// rendered badge, so this locates the badge URL and reads the percentage from
// it. It returns an empty string for non-Dependabot pull requests or when the
// score cannot be determined.
func CompatibilityScore(ctx context.Context, body string) string {
	badge := compatBadgeExpr.FindString(body)
	if badge == "" {
		return ""
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, badge, nil) //nolint:gosec // G107: URL is matched against the fixed dependabot-badges host
	if err != nil {
		return ""
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	data, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))

	match := compatScoreExpr.FindSubmatch(data)
	if match == nil {
		return ""
	}

	return string(match[1]) + "%"
}
