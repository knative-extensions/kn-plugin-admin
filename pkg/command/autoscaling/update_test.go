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
	"context"
	"testing"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/kn-plugin-admin/pkg"

	"knative.dev/kn-plugin-admin/pkg/testutil"
)

func TestNewAsUpdateSetCommand(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configAutoscaler,
			Namespace: knativeServing,
		},
		Data: make(map[string]string),
	}

	t.Run("no flags", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)

		_, err := testutil.ExecuteCommand(cmd)
		assert.ErrorContains(t, err, "'autoscaling update' requires flag(s)", err)
	})

	t.Run("operator mode should not be supported", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(cm)

		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodOperator,
		}
		cmd := NewAutoscalingUpdateCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.ErrorContains(t, err, "Knative managed by operator is not supported yet", err)
	})

	t.Run("config map not exist", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.ErrorContains(t, err, "failed to get ConfigMaps", err)
	})

	t.Run("enable scale-to-zero successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"enable-scale-to-zero": "false",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["enable-scale-to-zero"]
		assert.Check(t, ok, "key %q should exists", "enable-scale-to-zero")
		assert.Equal(t, "true", v, "enable-scale-to-zero should be true")
	})

	t.Run("disable scale-to-zero successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"enable-scale-to-zero": "true",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--no-scale-to-zero")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["enable-scale-to-zero"]
		assert.Check(t, ok, "key %q should exists", "enable-scale-to-zero")
		assert.Equal(t, "false", v, "enable-scale-to-zero should be false")
	})

	t.Run("enable scale-to-zero but it's already enabled", func(t *testing.T) {
		cm.Data = map[string]string{
			"enable-scale-to-zero": "true",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.NilError(t, err)

		updated, err := client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, equality.Semantic.DeepEqual(updated, cm), "configmap should not be changed")

	})

	t.Run("update container-concurrency-target-percentage successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"container-concurrency-target-percentage": "0.5",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--container-concurrency-target-percentage", "0.7")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["container-concurrency-target-percentage"]
		assert.Check(t, ok, "key %q should exists", "container-concurrency-target-percentage")
		assert.Equal(t, "0.7", v, "container-concurrency-target-percentage should be 0.7")
	})

	t.Run("update stable-window successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"stable-window": "60",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--stable-window", "2m")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["stable-window"]
		assert.Check(t, ok, "key %q should exists", "stable-window")
		assert.Equal(t, "120s", v, "stable-window should be 120s")
	})

	t.Run("return error if set stable-window less than 6s", func(t *testing.T) {
		cm.Data = map[string]string{
			"stable-window": "60",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--stable-window", "2s")
		assert.ErrorContains(t, err, "stable-window = 2s, must be in", err)
	})

	t.Run("return error if set max-scale-up-rate less than 1.0", func(t *testing.T) {
		cm.Data = map[string]string{
			"max-scale-up-rate": "2.0",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--max-scale-up-rate", "0.5")
		assert.ErrorContains(t, err, "max-scale-up-rate = 0.5, must be greater than 1.0", err)
	})

	t.Run("return error if set max-scale-down-rate less than 1.0", func(t *testing.T) {
		cm.Data = map[string]string{
			"max-scale-down-rate": "2.0",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--max-scale-up-rate", "0.5")
		assert.ErrorContains(t, err, "max-scale-up-rate = 0.5, must be greater than 1.0", err)
	})

	t.Run("return error if set scale-to-zero-grace-period less than 6s", func(t *testing.T) {
		cm.Data = map[string]string{
			"scale-to-zero-grace-period": "30s",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero-grace-period", "1s")
		assert.ErrorContains(t, err, "scale-to-zero-grace-period must be at least 6s, got 1s", err)
	})

	t.Run("return error if scale-to-zero-grace-period is not time duration", func(t *testing.T) {
		cm.Data = map[string]string{
			"scale-to-zero-grace-period": "30s",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero-grace-period", "60")
		assert.ErrorContains(t, err, "missing unit in duration 60", err)
	})

	t.Run("update scale-to-zero-pod-retention-period successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"scale-to-zero-pod-retention-period": "30s",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero-pod-retention-period", "1m")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["scale-to-zero-pod-retention-period"]
		assert.Check(t, ok, "key %q should exists", "scale-to-zero-pod-retention-period")
		assert.Equal(t, "60s", v, "scale-to-zero-pod-retention-period should be 60s")
	})

	t.Run("return error target-burst-capacity if set to -5", func(t *testing.T) {
		cm.Data = map[string]string{
			"target-burst-capacity": "-1",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--target-burst-capacity", "-5")
		assert.ErrorContains(t, err, "target-burst-capacity must be either non-negative or -1 (for unlimited), got -5", err)
	})

	t.Run("update pod-autoscaler-class successfully", func(t *testing.T) {
		cm.Data = map[string]string{
			"pod-autoscaler-class": "old.class",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--pod-autoscaler-class", "new.class")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["pod-autoscaler-class"]
		assert.Check(t, ok, "key %q should exists", "pod-autoscaler-class")
		assert.Equal(t, "new.class", v, "pod-autoscaler-classshould be new.class")
	})

	t.Run("return error if set activator-capacity to less than 1", func(t *testing.T) {
		cm.Data = map[string]string{
			"activator-capacity": "2",
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet:          client,
			InstallationMethod: pkg.InstallationMethodStandalone,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--activator-capacity", "0.5")
		assert.ErrorContains(t, err, "activator-capacity = 0.5, must be at least 1", err)
	})
}
