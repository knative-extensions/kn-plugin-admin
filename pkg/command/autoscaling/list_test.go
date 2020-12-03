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
	"sort"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client/pkg/util"
	"knative.dev/kn-plugin-admin/pkg"
	"knative.dev/kn-plugin-admin/pkg/testutil"
	"knative.dev/serving/pkg/autoscaler/config"
)

func checkListOutput(t *testing.T, data map[string]string, output string, noHeaders bool) {
	config, err := config.NewConfigFromMap(data)
	assert.NilError(t, err)

	lines := strings.Split(output, "\n")
	if !noHeaders {
		assert.Check(t, util.ContainsAll(lines[0], "NAME", "VALUE"))
	}

	names := make([]string, 0, len(configNameValueOfMap))
	for key := range configNameValueOfMap {
		names = append(names, key)
	}
	sort.Strings(names)

	start := 1
	if noHeaders {
		start = 0
	}
	for i, name := range names {
		value := configNameValueOfMap[name](config)
		assert.Check(t, util.ContainsAll(lines[i+start], name, value))
	}
}

func TestDescribesDuration(t *testing.T) {
	t0 := 123 * time.Second
	assert.Equal(t, t0.String(), describeDuration(t0))

	t1 := 60 * time.Second
	assert.Equal(t, "1m", describeDuration(t1))

	t2 := 3600 * time.Second
	assert.Equal(t, "1h", describeDuration(t2))

	t3 := 1238548 * time.Second
	assert.Equal(t, t3.String(), describeDuration(t3))
}

func TestAutoscalingListDefaultValues(t *testing.T) {
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
		checkListOutput(t, cm.Data, output, false)
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
		checkListOutput(t, cm.Data, output, false)
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
		checkListOutput(t, cm.Data, output, true)
	})
}
