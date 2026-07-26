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

// Package config provides common configuration types and helpers.
package config

import (
	"sync"
)

const (
	globalConfigFile = "git-repo.yml"
)

var (
	instance  *Configuration //nolint:gochecknoglobals // lazy-loaded configuration singleton
	loadError error          //nolint:gochecknoglobals // lazy-loaded configuration singleton
	mutex     sync.RWMutex   //nolint:gochecknoglobals // lazy-loaded configuration singleton
	once      sync.Once      //nolint:gochecknoglobals // lazy-loaded configuration singleton
)
