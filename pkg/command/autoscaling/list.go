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
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"knative.dev/client/pkg/kn/commands/flags"
	hprinters "knative.dev/client/pkg/printers"
	"knative.dev/kn-plugin-admin/pkg"
)

var (
	configNameList = []string{
		"activator-capacity",
		"container-concurrency-target-default",
		"enable-scale-to-zero",
		"max-scale-up-rate",
		"max-scale-down-rate",
		"panic-window-percentage",
		"panic-threshold-percentage",
		"pod-autoscaler-class",
		"requests-per-second-target-default",
		"stable-window",
		"scale-to-zero-grace-period",
		"scale-to-zero-pod-retention-period",
		"target-burst-capacity",
	}
)

// autoscalingListHandlers handles for `kn autoscaling list` command's output
func autoscalingListHandlers(h hprinters.PrintHandler) {
	autoscalingColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Description: "Name of the Autoscaling config", Priority: 1},
		{Name: "Value", Type: "string", Description: "Value of the Autoscaling config", Priority: 1},
	}
	h.TableHandler(autoscalingColumnDefinitions, printAutoscalingConfigs)
}

// printAutoscalingConfigs builds autoscaling config list table rows
func printAutoscalingConfigs(cm *corev1.ConfigMap, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(configNameList))
	for _, name := range configNameList {
		row := metav1beta1.TableRow{}
		if value := cm.Data[name]; value != "" {
			row.Cells = append(row.Cells, name, value)
			rows = append(rows, []metav1beta1.TableRow{row}...)
		}
	}
	return rows, nil
}

// NewAutoscalingListCommand represents autoscaling list command
func NewAutoscalingListCommand(p *pkg.AdminParams) *cobra.Command {
	autoscalingListFlags := flags.NewListPrintFlags(autoscalingListHandlers)
	autoscalingListCmd := &cobra.Command{
		Use:   "list",
		Short: "List autoscaling config",
		Long:  `List autoscaling config provided by Knative Pod Autoscaler (KPA)`,
		Example: `
  # To list all autoscaling configs 
  kn admin autoscaling list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get ConfigMaps: %+v", err)
			}

			err = autoscalingListFlags.Print(currentCm, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}

	autoscalingListFlags.HumanReadableFlags.AddFlags(autoscalingListCmd)
	autoscalingListFlags.GenericPrintFlags.OutputFlagSpecified = func() bool {
		return false
	}
	return autoscalingListCmd
}
