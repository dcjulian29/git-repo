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

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// useTempHome points the home-directory lookup at a fresh temporary directory
// so configuration reads and writes never touch the developer's real config.
func useTempHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir uses USERPROFILE on Windows

	return home
}

func TestLoadMissingFileReturnsEmptyConfig(t *testing.T) {
	useTempHome(t)

	cfg, err := loadFromDisk()
	if err != nil {
		t.Fatalf("unexpected error loading absent config: %v", err)
	}

	if cfg.Directory != "" || len(cfg.Repositories) != 0 {
		t.Fatalf("expected empty configuration, got %+v", cfg)
	}
}

func TestSaveThenLoadRoundTrip(t *testing.T) {
	home := useTempHome(t)

	if err := os.MkdirAll(filepath.Join(home, ".config"), 0o755); err != nil {
		t.Fatal(err)
	}

	want := &Configuration{
		Directory: "~/src",
		Repositories: []Repository{
			{Name: "alpha", URL: "https://example.com/alpha.git"},
			{Name: "beta", URL: "git@example.com:beta.git"},
		},
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := loadFromDisk()
	if err != nil {
		t.Fatalf("load after save failed: %v", err)
	}

	if got.Directory != want.Directory {
		t.Fatalf("directory = %q, want %q", got.Directory, want.Directory)
	}

	if len(got.Repositories) != len(want.Repositories) {
		t.Fatalf("got %d repositories, want %d", len(got.Repositories), len(want.Repositories))
	}

	for i, r := range want.Repositories {
		if got.Repositories[i] != r {
			t.Fatalf("repository %d = %+v, want %+v", i, got.Repositories[i], r)
		}
	}
}

func TestSaveNilConfigIsRejected(t *testing.T) {
	useTempHome(t)

	if err := Save(nil); err == nil {
		t.Fatal("expected an error when saving a nil configuration")
	}
}
