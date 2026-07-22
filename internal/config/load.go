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
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/yaml.v3"
)

// Load returns the configuration for the application.
// The configuration is loaded once and cached for subsequent calls.
// Returns an error if the file cannot be read or parsed.
func Load() (Configuration, error) {
	once.Do(func() {
		instance, loadError = load()
	})

	mutex.RLock()
	defer mutex.RUnlock()

	return *instance, loadError
}

func load() (*Configuration, error) {
	cfg := &Configuration{}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}

	filePath := filepath.Join(home, ".config", globalConfigFile)
	if !filesystem.FileExist(filePath) {
		return cfg, nil
	}

	file, err := os.ReadFile(filePath) //nolint:gosec // G304: reads the app's own config file (home dir + constant name), not external input

	if err != nil {
		return cfg, fmt.Errorf("could not read git repository configuration: %w", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to load git repository configuration: %w", err)
	}

	return cfg, nil
}
