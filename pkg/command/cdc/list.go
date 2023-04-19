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

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"knative.dev/client-pkg/pkg/kn/commands/flags"
	"knative.dev/client-pkg/pkg/printers"
	hprinters "knative.dev/client-pkg/pkg/printers"
	"knative.dev/kn-plugin-admin/pkg"
	typev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func cdcListHandlers(h printers.PrintHandler) {
	cdcListColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Domain Name", Type: "string", Description: "Name of the cluster domain claim object", Priority: 1},
		{Name: "Namespace", Type: "string", Description: "Namespace of the domain", Priority: 1},
	}
	h.TableHandler(cdcListColumnDefinitions, printCdcList)
}

func printCdcList(cdcList *typev1alpha1.ClusterDomainClaimList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(cdcList.Items))
	for _, item := range cdcList.Items {
		row := metav1beta1.TableRow{}
		row.Cells = append(row.Cells, item.Name, item.Spec.Namespace)
		rows = append(rows, row)
	}
	return rows, nil
}

// NewCdcListCommand represents 'kn-admin cdc list' command
func NewCdcListCommand(p *pkg.AdminParams) *cobra.Command {

	cdcListFlags := flags.NewListPrintFlags(cdcListHandlers)
	cdcListCommand := &cobra.Command{
		Use:   "list",
		Short: "List cluster domain claims",
		Long:  "List Knative cluster domain claims",
		Example: `
  # To list all cluster domain claims
  kn admin cdc list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := p.NewNetworkingClient()
			if err != nil {
				return err
			}
			cdcList, err := client.NetworkingV1alpha1().ClusterDomainClaims().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return err
			}
			err = cdcListFlags.Print(cdcList, cmd.OutOrStdout())
			return err
		},
	}
	cdcListFlags.HumanReadableFlags.AddFlags(cdcListCommand)
	cdcListFlags.GenericPrintFlags.OutputFlagSpecified = func() bool {
		return false
	}
	return cdcListCommand
}
