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
	"io/fs"
	"os"
	"path/filepath"
)

// FindGitRepositories walks path recursively and returns the absolute path of every
// directory that contains a ".git" sub-directory (i.e. every repository root).
// Unreadable entries are silently skipped.
func FindGitRepositories(path string) ([]string, error) {
	dirs := []string{}

	err := filepath.WalkDir(path, func(entry string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		if d.IsDir() && d.Name() == ".git" {
			dirs = append(dirs, filepath.Dir(entry))
			return filepath.SkipDir
		}

		return nil
	})

	return dirs, err
}
