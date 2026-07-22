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

package shared

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens rawURL in the user's default browser.
func OpenBrowser(rawURL string) error {
	name, args := browserCommand(runtime.GOOS, rawURL)

	return exec.Command(name, args...).Start() //nolint:gosec // G204: launcher is a fixed per-OS command; only the URL varies
}

// browserCommand returns the launcher and arguments used to open a URL on the
// given operating system.
func browserCommand(goos, rawURL string) (string, []string) {
	switch goos {
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", rawURL}
	case "darwin":
		return "open", []string{rawURL}
	default:
		return "xdg-open", []string{rawURL}
	}
}
