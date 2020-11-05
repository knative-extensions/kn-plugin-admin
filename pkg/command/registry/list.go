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

package registry

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"knative.dev/client/pkg/kn/commands"
	"knative.dev/client/pkg/kn/commands/flags"
	"knative.dev/kn-plugin-admin/pkg"
)

// NewRegistryListCommand represents the list command
func NewRegistryListCommand(p *pkg.AdminParams) *cobra.Command {
	registryListFlags := flags.NewListPrintFlags(RegistryListHandlers)
	var serviceaccount string
	var registryListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List registry settings",
		Long:    `List registry settings with server and username.`,
		Example: `
  # To list registry settings
  kn admin registry list \
    --namespace=[NAMESPACE] \
    --serviceaccount=[SERVICE_ACCOUNT]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// retrieve namespaces
			namespace := cmd.Flag("namespace").Value.String()
			if namespace == "" && serviceaccount != "" {
				return fmt.Errorf("cannot specifiy service account with empty namespace")
			}

			namespacesList, err := searchNamespace(p.ClientSet, namespace)
			if err != nil {
				return fmt.Errorf("failed to search specified namespaces: %v", err)
			}

			secretList := &corev1.SecretList{}
			for _, ns := range namespacesList.Items {
				err = addSecrets(p.ClientSet, ns.Name, serviceaccount, secretList)
			}

			// empty namespace indicates all-namespaces flag is specified
			if namespace == "" {
				registryListFlags.EnsureWithNamespace()
			}

			// Sort secretList by namespace and name (in this order)
			sort.SliceStable(secretList.Items, func(i, j int) bool {
				a := secretList.Items[i]
				b := secretList.Items[j]

				if a.Namespace != b.Namespace {
					return a.Namespace < b.Namespace
				}
				return a.ObjectMeta.Name < b.ObjectMeta.Name
			})

			return registryListFlags.Print(secretList, cmd.OutOrStdout())
		},
	}
	commands.AddNamespaceFlags(registryListCmd.Flags(), false)
	registryListFlags.HumanReadableFlags.AddFlags(registryListCmd)
	registryListFlags.GenericPrintFlags.OutputFlagSpecified = func() bool {
		return false
	}
	registryListCmd.Flags().StringVar(&serviceaccount, "serviceaccount", "", "the service account to save imagePullSecrets")
	registryListCmd.InitDefaultHelpFlag()
	return registryListCmd
}

func searchNamespace(kubeclient kubernetes.Interface, expectNamespace string) (*corev1.NamespaceList, error) {
	list, err := kubeclient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if expectNamespace != "" {
		for _, ns := range list.Items {
			if expectNamespace == ns.Name {
				return &corev1.NamespaceList{
					Items: []corev1.Namespace{ns},
				}, nil
			}
		}
		return nil, fmt.Errorf("namespace %s not found", expectNamespace)
	}
	return list, nil
}

func addSecrets(kubeclient kubernetes.Interface, ns string, sa string, secretList *corev1.SecretList) error {
	secrets, err := kubeclient.CoreV1().Secrets(ns).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(AdminRegistryLabels).String(),
	})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		if sa == "" || (sa == secret.Labels[ImagePullServiceAccount]) {
			secretList.Items = append(secretList.Items, secret)
		}
	}
	return nil
}
