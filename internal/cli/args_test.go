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

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestWithUsage(t *testing.T) {
	newCmd := func() (*cobra.Command, *bytes.Buffer) {
		cmd := &cobra.Command{Use: "widget <name>"}
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetErr(buf)

		return cmd, buf
	}

	t.Run("prints usage when validation fails", func(t *testing.T) {
		cmd, buf := newCmd()

		err := WithUsage(cobra.ExactArgs(1))(cmd, nil)
		if err == nil {
			t.Fatal("expected an error for missing argument")
		}

		if !strings.Contains(buf.String(), "Usage:") {
			t.Fatalf("expected usage output, got %q", buf.String())
		}
	})

	t.Run("stays quiet when validation passes", func(t *testing.T) {
		cmd, buf := newCmd()

		if err := WithUsage(cobra.ExactArgs(1))(cmd, []string{"a"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if buf.Len() != 0 {
			t.Fatalf("expected no output, got %q", buf.String())
		}
	})
}
