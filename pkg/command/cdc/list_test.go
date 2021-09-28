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
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/kn-plugin-admin/pkg/testutil"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func TestCdcListCommandWithoutKubeContext(t *testing.T) {
	t.Run("kubectl context is not set", func(t *testing.T) {
		p := testutil.NewTestAdminWithoutKubeConfig()
		cmd := NewCdcListCommand(p)
		_, err := testutil.ExecuteCommand(cmd)
		assert.Error(t, err, testutil.ErrNoKubeConfiguration)
	})
}

func TestCdcListSuccess(t *testing.T) {
	t.Run("list cdc", func(t *testing.T) {

		cdc := &v1alpha1.ClusterDomainClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: testDomain,
			},
			Spec: v1alpha1.ClusterDomainClaimSpec{
				Namespace: testNs,
			},
		}
		p := testutil.NewTestAdminParamsWithNetworkingObjects(cdc)
		cmd := NewCdcListCommand(p)
		out, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(out, testDomain))
		assert.Check(t, strings.Contains(out, testNs))
	})
}
