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
	"fmt"
	"strings"

	"gotest.tools/v3/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"

	"testing"

	"knative.dev/kn-plugin-admin/pkg/testutil"
)

func TestCdcDeleteCommandWithoutKubeContext(t *testing.T) {
	t.Run("kubectl context is not set", func(t *testing.T) {
		p := testutil.NewTestAdminWithoutKubeConfig()
		cmd := NewCdcDeleteCommand(p)
		_, err := testutil.ExecuteCommand(cmd)
		assert.Error(t, err, testutil.ErrNoKubeConfiguration)
	})
}

func TestCdcDeleteCommand(t *testing.T) {
	cdc := &v1alpha1.ClusterDomainClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: testDomain,
		},
		Spec: v1alpha1.ClusterDomainClaimSpec{
			Namespace: testNs,
		},
	}
	p := testutil.NewTestAdminParamsWithNetworkingObjects(cdc)
	t.Run("delete cdc successfully", func(t *testing.T) {
		cmd := NewCdcDeleteCommand(p)
		out, err := testutil.ExecuteCommand(cmd, testDomain)
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(out, fmt.Sprintf("'%s' deleted", testDomain)))
	})
	t.Run("delete cdc which does not exist", func(t *testing.T) {
		cmd := NewCdcDeleteCommand(p)
		notFoundDomain := "notfound.com"
		_, err := testutil.ExecuteCommand(cmd, notFoundDomain)
		assert.ErrorType(t, err, errors.IsNotFound)
	})
	t.Run("incomplete arg for cdc delete", func(t *testing.T) {
		p, _ := testutil.NewTestAdminParams()
		cmd := NewCdcDeleteCommand(p)
		_, err := testutil.ExecuteCommand(cmd)
		assert.ErrorContains(t, err, "cdc delete", "single argument")
		_, err = testutil.ExecuteCommand(cmd, "abc.com", "xyz.com")
		assert.ErrorContains(t, err, "cdc delete", "single argument")
	})
}
