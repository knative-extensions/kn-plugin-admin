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
	"knative.dev/kn-plugin-admin/pkg/testutil"
)

func TestNewCdcCreateCommand(t *testing.T) {
	name := "test.com"
	namespace := "test-ns"

	t.Run("kubectl context is not set", func(t *testing.T) {
		p := testutil.NewTestAdminWithoutKubeConfig()
		cmd := NewCdcCreateCommand(p)
		_, err := testutil.ExecuteCommand(cmd, name, "--namespace", namespace)
		assert.ErrorContains(t, err, testutil.ErrNoKubeConfiguration)
	})
	t.Run("incomplete args for cdc create", func(t *testing.T) {
		p, _ := testutil.NewTestAdminParams()
		cmd := NewCdcCreateCommand(p)
		_, err := testutil.ExecuteCommand(cmd, name)
		assert.ErrorContains(t, err, "required flag", "namespace")
	})
	t.Run("successful cdc create", func(t *testing.T) {
	})
}
