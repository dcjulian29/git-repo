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

// Config is the top-level structure that contains the root directory and
// the repositories managed within that directory
type Config struct {
	// Directory is the root directory that contains (or will contain) all
	// managed repositories.
	Directory string `yaml:"directory"`

	// Repositories is the list of repositories that are managed by this tool.
	Repositories []Repository `yaml:"repositories"`
}
