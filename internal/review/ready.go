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

package review

import (
	"context"
	"fmt"

	"github.com/fatih/color"
)

// ReadyPull marks a draft pull request as ready for review. When the pull
// request is already ready it reports so and makes no change.
func ReadyPull(ctx context.Context, ref Ref) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	pull, err := client.GetPull(ctx, ref.Repo, ref.Number)
	if err != nil {
		return err
	}

	if !pull.Draft {
		fmt.Printf("%s is already ready for review.\n", color.CyanString("%s#%d", ref.Name, ref.Number))

		return nil
	}

	if err := client.MarkPullReady(ctx, pull.NodeID); err != nil {
		return err
	}

	fmt.Println(color.GreenString("Marked %s#%d ready for review.", ref.Name, ref.Number))

	return nil
}
