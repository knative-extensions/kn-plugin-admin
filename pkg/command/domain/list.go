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

package domain

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/client-pkg/pkg/commands/flags"
	"knative.dev/kn-plugin-admin/pkg"
)

// NewDomainListCommand represents 'kn-admin domain list' command
func NewDomainListCommand(p *pkg.AdminParams) *cobra.Command {

	domainListFlags := flags.NewListPrintFlags(DomainListHandlers)
	domainListCommand := &cobra.Command{
		Use:   "list",
		Short: "List domain",
		Long:  "List Knative custom domain",
		Example: `
  # To list all custom domains
  kn admin domain list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := p.NewKubeClient()
			if err != nil {
				return err
			}

			domainCm, err := client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configDomain, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get ConfigMap %s in namespace %s: %+v", configDomain, knativeServing, err)
			}
			domainCmType := metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			}
			domainCm.TypeMeta = domainCmType
			err = domainListFlags.Print(domainCm, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	domainListFlags.HumanReadableFlags.AddFlags(domainListCommand)
	domainListFlags.GenericPrintFlags.OutputFlagSpecified = func() bool {
		return false
	}
	return domainListCommand
}
