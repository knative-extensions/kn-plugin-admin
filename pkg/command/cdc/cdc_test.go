// Copyright 2021 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cdc

import (
	"testing"

	"gotest.tools/v3/assert"
)

const (
	testDomain = "test.com"
	testNs     = "test-ns"
)

func TestNewCdcCmd(t *testing.T) {
	cmd := NewCdcCommand(nil)
	assert.Check(t, cmd.HasSubCommands(), "cmd cdc should have subcommands")
	assert.Equal(t, 3, len(cmd.Commands()), "cdc command should have 3 subcommands")

	_, _, err := cmd.Find([]string{"create"})
	assert.NilError(t, err, "cdc command should have create subcommand")

	_, _, err = cmd.Find([]string{"delete"})
	assert.NilError(t, err, "cdc command should have delete subcommand")

	_, _, err = cmd.Find([]string{"list"})
	assert.NilError(t, err, "cdc command should have list subcommand")
}
