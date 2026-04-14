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
	"errors"
	"os"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/yaml.v3"
)

// Save persists the provided configuration and updates the cached instance.
// Returns an error if configuration is nil or if the file cannot be written.
func Save(cfg *Configuration) error {
	if cfg == nil {
		return errors.New("can not save an uninitialized configuration")
	}

	mutex.Lock()
	defer mutex.Unlock()

	yaml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(home, ".config", globalConfigFile)

	if err := filesystem.EnsureFileExist(filePath, yaml); err != nil {
		return err
	}

	instance = cfg

	return nil
}
