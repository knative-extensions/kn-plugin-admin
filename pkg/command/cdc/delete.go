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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/kn-plugin-admin/pkg"
)

// NewCdcDeleteCommand to delete a cluster domain claims
func NewCdcDeleteCommand(p *pkg.AdminParams) *cobra.Command {
	cdcDeleteCommand := &cobra.Command{
		Use:   "delete",
		Short: "delete cluster domain claim",
		Long:  "Delete Knative cluster domain claim",
		Example: `
  # To delete a cluster domain claim
  kn admin cdc delete domain.name`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := p.NewNetworkingClient()
			if err != nil {
				return err
			}

			if len(args) != 1 {
				return errors.New("'cdc delete' requires the cdc name given as single argument")
			}
			name := args[0]
			err = client.NetworkingV1alpha1().ClusterDomainClaims().Delete(context.TODO(), name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Cluster Domain Claim '%s' deleted.\n", name)
			return nil
		},
	}
	return cdcDeleteCommand
}
