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
	"path/filepath"
	"testing"
)

func TestHasUpstream(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	if !HasUpstream(dir) {
		t.Fatal("expected tracking branch to report an upstream")
	}

	runGit(t, dir, "switch", "-c", "feature/no-upstream")

	if HasUpstream(dir) {
		t.Fatal("expected branch with no tracking branch to report no upstream")
	}
}

func TestComputeStatusCleanInSync(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	got := ComputeStatus(dir)

	want := RepoStatus{Folder: dir}
	want.Folder = got.Folder // path is echoed back verbatim

	if got != want {
		t.Fatalf("clean in-sync repo should have no flags set, got %+v", got)
	}
}

func TestComputeStatusDirty(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	writeFile(t, dir, "README.md", "changed\n")

	got := ComputeStatus(dir)

	if !got.Dirty {
		t.Fatalf("expected Dirty, got %+v", got)
	}
}

func TestComputeStatusUntracked(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	writeFile(t, dir, "new.txt", "untracked\n")

	got := ComputeStatus(dir)

	if !got.Untracked {
		t.Fatalf("expected Untracked, got %+v", got)
	}

	if got.Dirty {
		t.Fatalf("an untracked-only file should not mark the tree Dirty, got %+v", got)
	}
}

func TestComputeStatusPushNeeded(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	commit(t, dir, "local.txt", "ahead\n")

	got := ComputeStatus(dir)

	if !got.PushNeeded || got.PullNeeded || got.Diverged {
		t.Fatalf("expected PushNeeded only, got %+v", got)
	}
}

func TestComputeStatusPullNeeded(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	first := newSyncedClone(t, remote)

	// A second clone pushes a new commit, leaving the first behind.
	second := newClone(t, remote)
	commit(t, second, "remote.txt", "from second\n")
	runGit(t, second, "push", "origin", "main")

	got := ComputeStatus(first)

	if !got.PullNeeded || got.PushNeeded || got.Diverged {
		t.Fatalf("expected PullNeeded only, got %+v", got)
	}
}

func TestComputeStatusDiverged(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	first := newSyncedClone(t, remote)

	// The first clone commits locally without pushing.
	commit(t, first, "local.txt", "local work\n")

	// A second clone pushes a different commit to the shared remote.
	second := newClone(t, remote)
	commit(t, second, "remote.txt", "remote work\n")
	runGit(t, second, "push", "origin", "main")

	got := ComputeStatus(first)

	if !got.Diverged || got.PushNeeded || got.PullNeeded {
		t.Fatalf("expected Diverged only, got %+v", got)
	}
}

func TestComputeStatusNoUpstream(t *testing.T) {
	requireGit(t)

	remote := newRemote(t)
	dir := newSyncedClone(t, remote)

	runGit(t, dir, "switch", "-c", "feature/no-upstream")

	got := ComputeStatus(dir)

	if !got.NoUpstream {
		t.Fatalf("expected NoUpstream, got %+v", got)
	}

	if got.PushNeeded || got.PullNeeded || got.Diverged {
		t.Fatalf("no-upstream branch must not set push/pull/diverged, got %+v", got)
	}
}

func TestFindGitRepositories(t *testing.T) {
	requireGit(t)

	root := t.TempDir()

	// Two real repositories, one nested under a sub-directory.
	runGit(t, root, "init", "-b", "main", filepath.Join(root, "repoA"))
	if err := os.MkdirAll(filepath.Join(root, "group"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "init", "-b", "main", filepath.Join(root, "group", "repoB"))

	// A plain directory that is not a repository.
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatal(err)
	}

	found, err := FindGitRepositories(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]bool{
		filepath.Join(root, "repoA"):            true,
		filepath.Join(root, "group", "repoB"): true,
	}

	if len(found) != len(want) {
		t.Fatalf("expected %d repositories, got %d: %v", len(want), len(found), found)
	}

	for _, f := range found {
		if !want[f] {
			t.Fatalf("unexpected repository discovered: %s", f)
		}
	}
}
