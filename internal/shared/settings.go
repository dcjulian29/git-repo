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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/yaml.v3"
)

var (
	instance *Config
	loadErr  error
	mutex    sync.RWMutex
	once     sync.Once
)

// GetSettings loads and returns the configuration for the managed git repositories
// The settings are loaded once and cached for subsequent calls. Returns an error
// if the file cannot be read or parsed.
func GetSettings() (Config, error) {
	once.Do(func() {
		instance, loadErr = load()
	})

	mutex.RLock()
	defer mutex.RUnlock()

	if instance == nil {
		return Config{}, loadErr
	}

	return *instance, loadErr
}

// SaveSettings persists the provided configuration and updates the cached instance.
// Returns an error if configuration is nil or if the file cannot be written.
func SaveSettings(cfg *Config) error {
	if cfg == nil {
		return errors.New("can not save settings with an uninitialized configuration")
	}

	mutex.Lock()
	defer mutex.Unlock()

	if err := save(cfg); err != nil {
		return err
	}

	instance = cfg

	return nil
}

func load() (*Config, error) {
	cfg := &Config{}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}

	filePath := filepath.Join(home, ".config", "git-repo.yml")
	if !filesystem.FileExist(filePath) {
		return cfg, nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, fmt.Errorf("could not read git repository configuration: %w", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to load git repository configuration: %w", err)
	}

	return cfg, nil
}

func save(cfg *Config) error {
	yaml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(home, ".config", "git-repo.yml")

	return filesystem.EnsureFileExist(filePath, yaml)
}
