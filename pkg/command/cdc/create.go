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
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/kn-plugin-admin/pkg"
	typev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/networking/pkg/client/clientset/versioned/typed/networking/v1alpha1"
)

// NewCdcCreateCommand to create cluster domain claims
func NewCdcCreateCommand(p *pkg.AdminParams) *cobra.Command {
	var namespace string
	cdcCreateCommand := &cobra.Command{
		Use:   "create",
		Short: "create cluster domain claim",
		Long:  "Create Knative cluster domain claim",
		Example: `
  # To create a cluster domain claim
  kn admin cdc create domain.name --namespace ns1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			restConfig, err := p.RestConfig()
			if err != nil {
				return err
			}
			client, err := v1alpha1.NewForConfig(restConfig)
			if err != nil {
				return err
			}

			name := args[0]
			cdc := typev1alpha1.ClusterDomainClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: typev1alpha1.ClusterDomainClaimSpec{
					Namespace: namespace,
				},
			}
			_, err = client.ClusterDomainClaims().Create(context.TODO(), &cdc, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Cluster Domain Claim '%s' created.\n", name)
			return nil
		},
	}
	cdcCreateCommand.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace which is allowed to create a DomainMapping using this ClusterDomainClaim's name.")
	cdcCreateCommand.MarkFlagRequired("namespace")
	return cdcCreateCommand
}
