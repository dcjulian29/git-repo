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

// git-repo is a CLI tool for managing multiple local Git repositories.
//
// It combines repository status inspection, synchronization, initialization
// from a YAML configuration into a single binary.
package main

import (
	"fmt"
	"os"

	"github.com/dcjulian29/git-repo/cmd"
	"github.com/dcjulian29/go-toolbox/textformat"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "\n"+textformat.Fatal(err.Error()))
		os.Exit(1)
	}
}
