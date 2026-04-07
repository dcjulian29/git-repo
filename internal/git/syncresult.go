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

// SyncResult holds the output produced while synchronizing a single
// repository (fetch, pull, push combined).
type SyncResult struct {
	// Folder is the absolute path to the repository root.
	Folder string

	// Output is the combined stdout/stderr of the sync operations, trimmed.
	Output string
}
