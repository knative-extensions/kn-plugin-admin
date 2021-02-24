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

package registry

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	"knative.dev/client/pkg/util"
	"knative.dev/kn-plugin-admin/pkg"
	"knative.dev/kn-plugin-admin/pkg/testutil"
)

var (
	fakeUsername1 = "fakeUsername1"
	fakeUsername2 = "fakeUsername2"
	fakeUsername3 = "fakeUsername3"
	fakeUsername4 = "fakeUsername4"

	fakePassword = "fakePassword"

	fakeEmail1 = "fakeEmail1"
	fakeEmail2 = "fakeEmail2"
	fakeEmail3 = "fakeEmail3"
	fakeEmail4 = "fakeEmail4"

	fakeServer1 = "fakeServer1"
	fakeServer2 = "fakeServer2"
	fakeServer3 = "fakeServer3"
	fakeServer4 = "fakeServer4"

	fakeSecretName1 = "fakeSecret1"
	fakeSecretName2 = "fakeSecret2"
	fakeSecretName3 = "fakeSecret3"
	fakeSecretName4 = "fakeSecret4"

	fakeServiceAccount = "fakeServiceAccount"

	fakeNamespace = "fakeNamespace"

	defaultServiceAccount = "default"
	defaultNamespace      = "default"
)

func TestNewRegistryListCommand(t *testing.T) {
	t.Run("kubectl context is not set", func(t *testing.T) {
		p := testutil.NewTestAdminWithoutKubeConfig()
		cmd := NewRegistryListCommand(p)
		_, err := testutil.ExecuteCommand(cmd, "--serviceaccount", fakeServiceAccount, "--namespace", fakeNamespace)
		assert.Error(t, err, testutil.ErrNoKubeConfiguration)
	})

	t.Run("list registries with service account only but no namespace specified", func(t *testing.T) {
		p, client := testutil.NewTestAdminParams()
		assert.Check(t, client != nil)
		cmd := NewRegistryListCommand(p)

		_, err := testutil.ExecuteCommand(cmd, "--serviceaccount", "fakeServiceAccount")
		assert.ErrorContains(t, err, "cannot specifiy service account with empty namespace")
	})

	t.Run("list registries with service account and namespace specified", func(t *testing.T) {
		client := fakeRegistry()

		p := &pkg.AdminParams{
			NewKubeClient: func() (kubernetes.Interface, error) {
				return client, nil
			},
		}

		cmd := NewRegistryListCommand(p)

		output, err := testutil.ExecuteCommand(cmd, "--serviceaccount", fakeServiceAccount, "--namespace", fakeNamespace)
		assert.NilError(t, err)
		outputRows := strings.Split(output, "\n")

		assert.Equal(t, len(outputRows), 3)
		assert.Check(t, util.ContainsAll(outputRows[0], "SERVICEACCOUNT", "SECRET", "USERNAME", "SERVER", "EMAIL"))
		assert.Check(t, util.ContainsAll(outputRows[1], fakeServiceAccount, fakeSecretName1, fakeUsername1, fakeServer1, fakeEmail1))
	})

	t.Run("list registries with namespace specified only", func(t *testing.T) {
		client := fakeRegistry()

		p := &pkg.AdminParams{
			NewKubeClient: func() (kubernetes.Interface, error) {
				return client, nil
			},
		}

		cmd := NewRegistryListCommand(p)

		output, err := testutil.ExecuteCommand(cmd, "--namespace", fakeNamespace)
		assert.NilError(t, err)
		outputRows := strings.Split(output, "\n")

		assert.Equal(t, len(outputRows), 4)
		assert.Check(t, util.ContainsAll(outputRows[0], "SERVICEACCOUNT", "SECRET", "USERNAME", "SERVER", "EMAIL"))
		assert.Check(t, util.ContainsAll(outputRows[1], fakeServiceAccount, fakeSecretName1, fakeUsername1, fakeServer1, fakeEmail1))
		assert.Check(t, util.ContainsAll(outputRows[2], defaultServiceAccount, fakeSecretName2, fakeUsername2, fakeServer2, fakeEmail2))
	})

	t.Run("list all registries in all namespaces and service accounts", func(t *testing.T) {
		client := fakeRegistry()

		p := &pkg.AdminParams{
			NewKubeClient: func() (kubernetes.Interface, error) {
				return client, nil
			},
		}

		cmd := NewRegistryListCommand(p)

		output, err := testutil.ExecuteCommand(cmd)
		assert.NilError(t, err)
		outputRows := strings.Split(output, "\n")

		assert.Equal(t, len(outputRows), 6)
		assert.Check(t, util.ContainsAll(outputRows[0], "NAMESPACE", "SERVICEACCOUNT", "SECRET", "USERNAME", "SERVER", "EMAIL"))
		// sorted by namespace
		assert.Check(t, util.ContainsAll(outputRows[3], fakeNamespace, fakeServiceAccount, fakeSecretName1, fakeUsername1, fakeServer1, fakeEmail1))
		assert.Check(t, util.ContainsAll(outputRows[4], fakeNamespace, defaultServiceAccount, fakeSecretName2, fakeUsername2, fakeServer2, fakeEmail2))
		assert.Check(t, util.ContainsAll(outputRows[1], defaultNamespace, defaultServiceAccount, fakeSecretName3, fakeUsername3, fakeServer3, fakeEmail3))
		assert.Check(t, util.ContainsAll(outputRows[2], defaultNamespace, defaultServiceAccount, fakeSecretName4, fakeUsername4, fakeServer4, fakeEmail4))
	})

}

func fakeRegistry() kubernetes.Interface {
	// secret1 listed by namespace and service account specified
	secret1 := createMockSecretWithParams(fakeSecretName1, fakeNamespace, fakeServiceAccount, fakeUsername1, fakeServer1, fakeEmail1)
	// secret1 and secret2 listed by namespace specified
	secret2 := createMockSecretWithParams(fakeSecretName2, fakeNamespace, defaultServiceAccount, fakeUsername2, fakeServer2, fakeEmail2)
	// secret1, secret2, secret3 and secret4 listed by no specifications
	// secret3 and secret4 exist in a same service account
	secret3 := createMockSecretWithParams(fakeSecretName3, defaultNamespace, defaultServiceAccount, fakeUsername3, fakeServer3, fakeEmail3)
	secret4 := createMockSecretWithParams(fakeSecretName4, defaultNamespace, defaultServiceAccount, fakeUsername4, fakeServer4, fakeEmail4)

	sa1 := createMockServiceAccountWithParams(
		fakeServiceAccount,
		fakeNamespace,
		[]corev1.LocalObjectReference{
			{
				Name: fakeSecretName1,
			},
		},
	)

	sa2 := createMockServiceAccountWithParams(
		defaultServiceAccount,
		fakeNamespace,
		[]corev1.LocalObjectReference{
			{
				Name: fakeSecretName2,
			},
		},
	)

	sa3 := createMockServiceAccountWithParams(
		defaultServiceAccount,
		defaultNamespace,
		[]corev1.LocalObjectReference{
			{
				Name: fakeSecretName3,
			},
			{
				Name: fakeSecretName4,
			},
		},
	)

	fakeNS := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fakeNamespace,
		},
	}
	defaultNS := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultNamespace,
		},
	}

	client := k8sfake.NewSimpleClientset(&defaultNS, &fakeNS, &sa1, &sa2, &sa3, &secret1, &secret2, &secret3, &secret4)

	return client
}

func createMockServiceAccountWithParams(name, namespace string, secrets []corev1.LocalObjectReference) corev1.ServiceAccount {
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		ImagePullSecrets: secrets,
	}
	return sa
}

func createMockSecretWithParams(name, namespace, serviceaccount, username, server, email string) corev1.Secret {
	dockerCfg := Registry{
		Auths: Auths{
			server: registryCred{
				Username: username,
				Password: "password",
				Email:    email,
			},
		},
	}

	j, _ := json.Marshal(dockerCfg)

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				pkg.LabelManagedBy:      AdminRegistryCmdName,
				ImagePullServiceAccount: serviceaccount,
			},
		},
		Data: map[string][]byte{
			".dockerconfigjson": j,
		},
	}
	return secret
}
