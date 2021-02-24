// Copyright 2020 The Knative Authors
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

package domain

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/client/pkg/util"

	"knative.dev/kn-plugin-admin/pkg/testutil"
)

func TestDomainListWithoutKubeContext(t *testing.T) {
	t.Run("kubectl context is not set", func(t *testing.T) {
		p := testutil.NewTestAdminWithoutKubeConfig()
		cmd := NewDomainListCommand(p)
		_, err := testutil.ExecuteCommand(cmd)
		assert.Error(t, err, testutil.ErrNoKubeConfiguration)
	})
}

func TestDomainListEmpty(t *testing.T) {
	t.Run("list domain", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{},
		}
		p, client := testutil.NewTestAdminParams(cm)
		assert.Check(t, client != nil)
		cmd := NewDomainListCommand(p)
		output, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		rowsOfOutput := strings.Split(output, "\n")
		assert.Check(t, util.ContainsAll(rowsOfOutput[0], "CUSTOM-DOMAIN", "SELECTOR"))
	})
}

func TestDomainListCommand(t *testing.T) {

	t.Run("list domain", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"test1.domain":  "",
				"a-test.domain": "selector:\n  app1: helloworld1\n app2: helloworld2\n",
				"test2.domain":  "selector:\n  app: helloworld\n",
			},
		}
		p, client := testutil.NewTestAdminParams(cm)
		assert.Check(t, client != nil)
		cmd := NewDomainListCommand(p)
		output, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		rowsOfOutput := strings.Split(output, "\n")
		//Domain will be listed with order by domain name
		assert.Check(t, util.ContainsAll(rowsOfOutput[0], "CUSTOM-DOMAIN", "SELECTOR"))
		assert.Check(t, util.ContainsAll(rowsOfOutput[1], "a-test.domain", "app1=helloworld1; app2=helloworld2"))
		assert.Check(t, util.ContainsAll(rowsOfOutput[2], "test1.domain"))
		assert.Check(t, util.ContainsAll(rowsOfOutput[3], "test2.domain", "app=helloworld"))
	})
}

func TestDomainListCommandNoHeader(t *testing.T) {

	t.Run("list domain", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"test1.domain": "",
				"test2.domain": "selector:\n  app: helloworld\n",
			},
		}
		p, client := testutil.NewTestAdminParams(cm)
		assert.Check(t, client != nil)
		cmd := NewDomainListCommand(p)
		output, err := testutil.ExecuteCommand(cmd, "--no-headers")
		assert.NilError(t, err)
		rowsOfOutput := strings.Split(output, "\n")
		assert.Check(t, util.ContainsAll(rowsOfOutput[0], "test1.domain"))
		assert.Check(t, util.ContainsAll(rowsOfOutput[1], "test2.domain", "app=helloworld"))
	})
}
