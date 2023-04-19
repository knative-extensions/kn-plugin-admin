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

	corev1 "k8s.io/api/core/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	hprinters "knative.dev/client-pkg/pkg/printers"
)

// RegistryListHandlers adds print handlers for registry list command
func RegistryListHandlers(h hprinters.PrintHandler) {
	registryColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Namespace", Type: "string", Description: "Namespace of the Knative service.", Priority: 0},
		{Name: "ServiceAccount", Type: "string", Description: "The ServiceAccount to save ImagePullSecrets.", Priority: 1},
		{Name: "Secret", Type: "string", Description: "The Secret to save registry.", Priority: 1},
		{Name: "UserName", Type: "string", Description: "The username of the registry.", Priority: 1},
		{Name: "Server", Type: "string", Description: "The server url of the registry.", Priority: 1},
		{Name: "Email", Type: "string", Description: "The email of the registry user.", Priority: 1},
	}

	h.TableHandler(registryColumnDefinitions, printRegistry)
	h.TableHandler(registryColumnDefinitions, printRegistryList)
}

// Private functions

// printRegistryList populates the registry list table rows
func printRegistryList(secretList *corev1.SecretList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(secretList.Items))

	for _, secret := range secretList.Items {
		r, err := printRegistry(&secret, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printRegistry populates the registry table rows
func printRegistry(secret *corev1.Secret, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: secret},
	}

	if options.AllNamespaces {
		row.Cells = append(row.Cells, secret.Namespace)
	}

	sa := secret.Labels[ImagePullServiceAccount]
	name := secret.Name
	registry := Registry{}
	err := json.Unmarshal(secret.Data[DockerJSONName], &registry)
	if err != nil {
		return []metav1beta1.TableRow{row}, err
	}
	for secretServer, secretAuth := range registry.Auths {
		row.Cells = append(row.Cells,
			sa,
			name,
			secretAuth.Username,
			secretServer,
			secretAuth.Email,
		)
	}

	return []metav1beta1.TableRow{row}, nil
}
