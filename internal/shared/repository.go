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

// Repository represents a single Git repository entry in the configuration
// file. Name is used as the local sub-directory name; URL is the remote that
// will be passed to "git clone".
type Repository struct {
	// Name is the local directory name for the cloned repository.
	Name string `yaml:"name"`

	// URL is the remote Git URL used when cloning.
	URL string `yaml:"url"`
}
