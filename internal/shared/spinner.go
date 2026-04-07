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
	"time"

	"github.com/theckman/yacspin"
)

// NewSpinner creates and returns a [yacspin.Spinner] configured with the
// project-standard style: 100 ms tick interval, yellow foreground, charset 69.
//
// Callers are responsible for calling spinner.Start() before work begins and
// defer spinner.Stop() (or an explicit Stop) when work is complete.
func NewSpinner() (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency: 100 * time.Millisecond,
		Colors:    []string{"fgYellow"},
		CharSet:   yacspin.CharSets[69],
	}

	return yacspin.New(cfg)
}
