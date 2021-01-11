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

package autoscaling

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"knative.dev/client/pkg/kn/commands/flags"
	hprinters "knative.dev/client/pkg/printers"
	"knative.dev/kn-plugin-admin/pkg"
	"knative.dev/serving/pkg/autoscaler/config"
	"knative.dev/serving/pkg/autoscaler/config/autoscalerconfig"
)

// A function for getting specific field value from autoscaler config
type valueOfConfig func(*autoscalerconfig.Config) string

var (
	ConfigNameValueOfMap = map[string]valueOfConfig{
		"activator-capacity": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.ActivatorCapacity)
		},
		"container-concurrency-target-default": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.ContainerConcurrencyTargetDefault)
		},
		"enable-scale-to-zero": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%+v", config.EnableScaleToZero)
		},
		"max-scale-up-rate": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.MaxScaleUpRate)
		},
		"max-scale-down-rate": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.MaxScaleDownRate)
		},
		"panic-window-percentage": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.PanicWindowPercentage)
		},
		"panic-threshold-percentage": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.PanicThresholdPercentage)
		},
		"pod-autoscaler-class": func(config *autoscalerconfig.Config) string {
			return config.PodAutoscalerClass
		},
		"requests-per-second-target-default": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.RPSTargetDefault)
		},
		"stable-window": func(config *autoscalerconfig.Config) string {
			return describeDuration(config.StableWindow)
		},
		"scale-to-zero-grace-period": func(config *autoscalerconfig.Config) string {
			return describeDuration(config.ScaleToZeroGracePeriod)
		},
		"scale-to-zero-pod-retention-period": func(config *autoscalerconfig.Config) string {
			return describeDuration(config.ScaleToZeroPodRetentionPeriod)
		},
		"target-burst-capacity": func(config *autoscalerconfig.Config) string {
			return fmt.Sprintf("%.1f", config.TargetBurstCapacity)
		},
	}
)

// describeDuration describes time.duration without 'm0s' and 'h0m'
func describeDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

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
	rows := make([]metav1beta1.TableRow, 0, len(ConfigNameValueOfMap))
	config, err := config.NewConfigFromMap(cm.Data)
	if err != nil {
		return rows, fmt.Errorf("failed to get autoscaling config: %+v", err)
	}

	// sort config names
	names := make([]string, 0, len(ConfigNameValueOfMap))
	for key := range ConfigNameValueOfMap {
		names = append(names, key)
	}
	sort.Strings(names)

	for _, name := range names {
		row := metav1beta1.TableRow{}
		row.Cells = append(row.Cells, name, ConfigNameValueOfMap[name](config))
		rows = append(rows, []metav1beta1.TableRow{row}...)
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
			client, err := p.NewKubeClient()
			if err != nil {
				return err
			}

			currentCm := &corev1.ConfigMap{}
			currentCm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
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
