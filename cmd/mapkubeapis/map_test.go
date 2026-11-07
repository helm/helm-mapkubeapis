/*
Copyright

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

package main

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestMapCmdIgnoresUnknownFlags(t *testing.T) {
	cmd := newMapCmd(io.Discard)

	var gotArgs []string
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		gotArgs = args
		return nil
	}
	cmd.SetArgs([]string{"my-release", "--namespace", "default", "--not-a-mapkubeapis-flag"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected unknown flags to be ignored, got error: %v", err)
	}
	if len(gotArgs) != 1 || gotArgs[0] != "my-release" {
		t.Fatalf("expected release name to be passed through, got %v", gotArgs)
	}
}
