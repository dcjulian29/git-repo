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

import "time"

// Item represents a single open issue or pull request.
type Item struct {
	// Number is the issue or pull-request number within its repository.
	Number int `json:"number"`

	// Title is the issue or pull-request title.
	Title string `json:"title"`

	// Author is the login of the user who opened it.
	Author string `json:"author"`

	// URL is the html_url that opens the item in a browser.
	URL string `json:"url"`

	// CreatedAt is the time the item was opened.
	CreatedAt time.Time `json:"created_at"`

	// IsPull is true when the item is a pull request rather than an issue.
	IsPull bool `json:"is_pull"`

	// Draft is true when the item is a draft pull request.
	Draft bool `json:"draft"`

	// Merged is true when the item is a pull request that has been merged.
	Merged bool `json:"merged"`

	// MergeableState is a pull request's coarse merge state ("clean",
	// "unstable", "dirty", ...). It is populated on demand for pull-request
	// listings and is empty for issues.
	MergeableState string `json:"mergeable_state,omitempty"`
}

// issuePayload mirrors the fields git-repo needs from the GitHub issues API,
// which returns both issues and pull requests. Pull requests carry a non-nil
// pull_request object.
type issuePayload struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	HTMLURL     string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"`
	Draft       bool      `json:"draft"`
	User        struct {
		Login string `json:"login"`
	} `json:"user"`
	PullRequest *struct {
		URL      string     `json:"url"`
		MergedAt *time.Time `json:"merged_at"`
	} `json:"pull_request"`
}

// toItem converts a decoded API payload into the exported Item type.
func (p issuePayload) toItem() Item {
	return Item{
		Number:    p.Number,
		Title:     p.Title,
		Author:    p.User.Login,
		URL:       p.HTMLURL,
		CreatedAt: p.CreatedAt,
		IsPull:    p.PullRequest != nil,
		Draft:     p.Draft,
		Merged:    p.PullRequest != nil && p.PullRequest.MergedAt != nil,
	}
}
