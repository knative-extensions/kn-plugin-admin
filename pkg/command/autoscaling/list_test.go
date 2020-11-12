// Copyright © 2020 The Knative Authors
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

package autoscaling

import (
	"strings"
	"testing"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client/pkg/util"
	"knative.dev/kn-plugin-admin/pkg"

	"knative.dev/kn-plugin-admin/pkg/testutil"
)

func TestAutoscalingListEmpty(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{ClientSet: client}
		cmd := NewAutoscalingListCommand(&p)
		output, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		lines := strings.Split(output, "\n")
		assert.Check(t, util.ContainsAll(lines[0], "NAME", "VALUE"))
	})
}

func TestAutoscalingListCommand(t *testing.T) {
	t.Run("list autoscaling configs", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"enable-scale-to-zero":    "true",
				"panic-window-percentage": "10",
				"max-scale-up-rate":       "100",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{ClientSet: client}
		cmd := NewAutoscalingListCommand(&p)
		output, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		lines := strings.Split(output, "\n")
		assert.Check(t, util.ContainsAll(lines[0], "NAME", "VALUE"))
		assert.Check(t, util.ContainsAll(lines[1], "enable-scale-to-zero", "true"))
		assert.Check(t, util.ContainsAll(lines[2], "max-scale-up-rate", "100"))
		assert.Check(t, util.ContainsAll(lines[3], "panic-window-percentage", "10"))
	})
}

func TestAutoscalingListCommandNoHeader(t *testing.T) {
	t.Run("list autoscaling configs", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"enable-scale-to-zero": "true",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{ClientSet: client}
		cmd := NewAutoscalingListCommand(&p)
		output, err := testutil.ExecuteCommand(cmd, "--no-headers")
		assert.NilError(t, err)
		lines := strings.Split(output, "\n")
		assert.Check(t, util.ContainsAll(lines[0], "enable-scale-to-zero", "true"))
	})
}
