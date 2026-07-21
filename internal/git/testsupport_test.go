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

package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// requireGit skips the calling test when the git executable is not available on
// PATH, so the integration tests degrade gracefully in minimal environments.
func requireGit(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found on PATH; skipping integration test")
	}
}

// runGit runs a git command in dir and fails the test on error. The environment
// is pinned to a deterministic identity and isolated from the developer's global
// and system git configuration so results do not depend on the host.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_CONFIG_SYSTEM="+os.DevNull,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s failed: %v\n%s", args, dir, err, out)
	}

	return string(out)
}

// newRemote creates a bare repository to act as an "origin" and returns its path.
func newRemote(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGit(t, dir, "init", "--bare", "-b", "main", ".")

	return dir
}

// newClone clones remote into a fresh temporary directory and returns the path
// to the working copy.
func newClone(t *testing.T, remote string) string {
	t.Helper()

	parent := t.TempDir()
	dir := filepath.Join(parent, "clone")
	runGit(t, parent, "clone", remote, dir)

	return dir
}

// writeFile creates or overwrites name inside dir with content.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

// commit writes name with content, stages it, and commits it in dir.
func commit(t *testing.T, dir, name, content string) {
	t.Helper()

	writeFile(t, dir, name, content)
	runGit(t, dir, "add", name)
	runGit(t, dir, "commit", "-m", "add "+name)
}

// newSyncedClone returns a clone that has one commit pushed to remote and its
// current branch tracking origin/main, i.e. a clean, in-sync repository.
func newSyncedClone(t *testing.T, remote string) string {
	t.Helper()

	dir := newClone(t, remote)
	commit(t, dir, "README.md", "hello\n")
	runGit(t, dir, "push", "-u", "origin", "main")

	return dir
}
