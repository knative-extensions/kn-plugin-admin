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

package pkg

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestAdminParams_installationMethod(t *testing.T) {
	t.Run("managed by KnativeOperator", func(t *testing.T) {
		ptrTrue := true
		domainCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         "operator.knative.dev/v1alpha1",
						Kind:               "KnativeServing",
						Name:               "knative-serving",
						BlockOwnerDeletion: &ptrTrue,
						Controller:         &ptrTrue,
					},
				},
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(domainCM)

		params := &AdminParams{
			NewKubeClient: func() (kubernetes.Interface, error) {
				return client, nil
			},
		}
		got, err := params.installationMethod()
		if err != nil {
			t.Error(err)
		}
		if got != InstallationMethodOperator {
			t.Error("Should return InstallationMethodOperator")
		}
	})

	t.Run("Installed by standalone mode", func(t *testing.T) {
		domainCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(domainCM)

		params := &AdminParams{
			NewKubeClient: func() (kubernetes.Interface, error) {
				return client, nil
			},
		}
		got, err := params.installationMethod()
		if err != nil {
			t.Error(err)
		}
		if got != InstallationMethodStandalone {
			t.Error("Should return InstallationMethodStandalone")
		}
	})

}

func TestAdminParams_EnsureInstallMethodStandalone(t *testing.T) {
	t.Run("Installation method unknown", func(t *testing.T) {
		params := &AdminParams{}
		if err := params.EnsureInstallMethodStandalone(); err == nil {
			t.Error("should return error for unknown installation method")
		}
	})

	t.Run("Installation method operator", func(t *testing.T) {
		params := &AdminParams{
			InstallationMethod: InstallationMethodOperator,
		}
		if err := params.EnsureInstallMethodStandalone(); err == nil {
			t.Error("should return error for operator installation")
		}
	})

	t.Run("Installation method standalone", func(t *testing.T) {
		params := &AdminParams{
			InstallationMethod: InstallationMethodStandalone,
		}
		if err := params.EnsureInstallMethodStandalone(); err != nil {
			t.Errorf("should not return error. got %#v", err)
		}
	})

}
