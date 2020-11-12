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
	"errors"
	"fmt"
	"strconv"
	"time"

	"knative.dev/kn-plugin-admin/pkg/command/utils"

	"knative.dev/kn-plugin-admin/pkg"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"knative.dev/client/pkg/kn/flags"

	as "knative.dev/serving/pkg/apis/autoscaling"
)

var (
	scaleToZero                                                         bool
	enableScaleToZero                                                   = "enable-scale-to-zero"
	knativeServing                                                      = "knative-serving"
	configAutoscaler                                                    = "config-autoscaler"
	stableWindow, scaleToZeroGracePeriod, scaleToZeroPodRetentionPeriod time.Duration

	containerConcurrencyTargetPercentage, panicWindowPercentage, panicThresholdPercentage, maxScaleUpRate,
	maxScaleDownRate, targetBurstCapacity, activatorCapacity, requestsPerSecondTargetDefault,
	containerConcurrencyTargetDefault, podAutoscalerClass string
)

func NewAutoscalingUpdateCommand(p *pkg.AdminParams) *cobra.Command {
	AutoscalingUpdateCommand := &cobra.Command{
		Use:   "update",
		Short: "Update autoscaling config",
		Long:  `Update autoscaling config provided by Knative Pod Autoscaler (KPA)`,
		Example: `
  # To enable scale-to-zero
  kn admin autoscaling update --scale-to-zero

  # To update stable window
  kn admin autoscaling update --stable-window 2m`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				return errors.New("'autoscaling update' requires flag(s)")
			}
			if err := p.EnsureInstallMethodStandalone(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get ConfigMaps: %+v", err)
			}
			desiredCm := currentCm.DeepCopy()

			if cmd.Flags().Changed("scale-to-zero") && cmd.Flags().Changed("no-scale-to-zero") {
				return fmt.Errorf("please specify either --scale-to-zero or --no-scale-to-zero")
			}

			if cmd.Flags().Changed("scale-to-zero") {
				desiredCm.Data[enableScaleToZero] = "true"
			}

			if cmd.Flags().Changed("no-scale-to-zero") {
				desiredCm.Data[enableScaleToZero] = "false"
			}

			if cmd.Flags().Changed("requests-per-second-target-default") {
				desiredCm.Data["requests-per-second-target-default"] = fmt.Sprintf("%s", requestsPerSecondTargetDefault)
			}

			if cmd.Flags().Changed("container-concurrency-target-default") {
				desiredCm.Data["container-concurrency-target-default"] = fmt.Sprintf("%s", containerConcurrencyTargetDefault)
			}

			if cmd.Flags().Changed("container-concurrency-target-percentage") {
				desiredCm.Data["container-concurrency-target-percentage"] = fmt.Sprintf("%s", containerConcurrencyTargetPercentage)
			}

			if cmd.Flags().Changed("stable-window") {
				if stableWindow < as.WindowMin || stableWindow > as.WindowMax {
					return fmt.Errorf("stable-window = %v, must be in [%v; %v] range", stableWindow,
						as.WindowMin, as.WindowMax)
				}

				if stableWindow.Round(time.Second) != stableWindow {
					return fmt.Errorf("stable-window = %v, must be specified with at most second precision", stableWindow)
				}

				fmt.Printf("debug stable-window %vs\n", stableWindow.Seconds())
				desiredCm.Data["stable-window"] = fmt.Sprintf("%vs", stableWindow.Seconds())
			}

			if cmd.Flags().Changed("panic-window-percentage") {
				desiredCm.Data["panic-window-percentage"] = panicWindowPercentage
			}

			if cmd.Flags().Changed("panic-threshold-percentage") {
				desiredCm.Data["panic-threshold-percentage"] = panicThresholdPercentage
			}

			if cmd.Flags().Changed("max-scale-up-rate") {
				tmp, err := strconv.ParseFloat(maxScaleUpRate, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", maxScaleUpRate)
				}
				if tmp <= 1.0 {
					return fmt.Errorf("max-scale-up-rate = %v, must be greater than 1.0", maxScaleUpRate)
				}
				desiredCm.Data["max-scale-up-rate"] = maxScaleUpRate
			}

			if cmd.Flags().Changed("max-scale-down-rate") {
				tmp, err := strconv.ParseFloat(maxScaleDownRate, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", maxScaleUpRate)
				}
				if tmp <= 1.0 {
					return fmt.Errorf("max-scale-down-rate = %v, must be greater than 1.0", maxScaleDownRate)
				}
				desiredCm.Data["max-scale-down-rate"] = maxScaleDownRate
			}

			if cmd.Flags().Changed("scale-to-zero-grace-period") {
				if scaleToZeroGracePeriod < as.WindowMin {
					return fmt.Errorf("scale-to-zero-grace-period must be at least %v, got %v", as.WindowMin, scaleToZeroGracePeriod)
				}

				desiredCm.Data["scale-to-zero-grace-period"] = fmt.Sprintf("%vs", scaleToZeroGracePeriod.Seconds())
			}

			if cmd.Flags().Changed("scale-to-zero-pod-retention-period") {
				if scaleToZeroPodRetentionPeriod < 0 {
					return fmt.Errorf("scale-to-zero-pod-retention-period cannot be negative, was: %v", scaleToZeroPodRetentionPeriod)
				}
				desiredCm.Data["scale-to-zero-pod-retention-period"] = fmt.Sprintf("%vs", scaleToZeroPodRetentionPeriod.Seconds())
			}

			if cmd.Flags().Changed("target-burst-capacity") {
				tmp, err := strconv.ParseFloat(targetBurstCapacity, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", maxScaleUpRate)
				}
				if tmp < 0 && tmp != -1 {
					return fmt.Errorf("target-burst-capacity must be either non-negative or -1 (for unlimited), got %s", targetBurstCapacity)
				}
				desiredCm.Data["target-burst-capacity"] = targetBurstCapacity
			}

			if cmd.Flags().Changed("pod-autoscaler-class") {
				desiredCm.Data["pod-autoscaler-class"] = podAutoscalerClass
			}

			if cmd.Flags().Changed("activator-capacity") {
				tmp, err := strconv.ParseFloat(activatorCapacity, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", maxScaleUpRate)
				}
				if tmp < 1 {
					return fmt.Errorf("activator-capacity = %v, must be at least 1", activatorCapacity)
				}
				desiredCm.Data["activator-capacity"] = activatorCapacity
			}

			err = utils.UpdateConfigMap(p.ClientSet, desiredCm)
			if err != nil {
				return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", configAutoscaler, knativeServing, err)
			}
			cmd.Printf("Updated Knative autoscaling config\n")

			return nil
		},
	}

	flags.AddBothBoolFlagsUnhidden(AutoscalingUpdateCommand.Flags(), &scaleToZero, "scale-to-zero", "", true,
		"Enable scale-to-zero if set.")
	AutoscalingUpdateCommand.Flags().StringVarP(&requestsPerSecondTargetDefault, "requests-per-second-target-default", "", "200", "the default target value for requests per second")
	AutoscalingUpdateCommand.Flags().StringVarP(&containerConcurrencyTargetDefault, "container-concurrency-target-default", "", "100", "the default value of container concurrency target")
	AutoscalingUpdateCommand.Flags().StringVarP(&containerConcurrencyTargetPercentage, "container-concurrency-target-percentage", "", "0.7", "percentage of the specified target should actually be targeted by the Autoscaler")
	AutoscalingUpdateCommand.Flags().DurationVarP(&stableWindow, "stable-window", "", 60*time.Second, "when operating in a stable mode, the autoscaler operates on the average concurrency over the x seconds of stable window")
	AutoscalingUpdateCommand.Flags().StringVarP(&panicWindowPercentage, "panic-window-percentage", "", "10", "The panic window is defined as a percentage of the stable window")
	AutoscalingUpdateCommand.Flags().StringVarP(&panicThresholdPercentage, "panic-threshold-percentage", "", "200", "This threshold defines when the autoscaler will move from stable mode into panic mode")
	AutoscalingUpdateCommand.Flags().StringVarP(&maxScaleUpRate, "max-scale-up-rate", "", "1000", "Maximum ratio of desired vs. observed pods")
	AutoscalingUpdateCommand.Flags().StringVarP(&maxScaleDownRate, "max-scale-down-rate", "", "2", "Maximum ratio of observed vs. desired pods")
	AutoscalingUpdateCommand.Flags().DurationVarP(&scaleToZeroGracePeriod, "scale-to-zero-grace-period", "", 30*time.Second, " the maximum seconds of time that the last pod will remain active after the Autoscaler has decided to scale pods to zero")
	AutoscalingUpdateCommand.Flags().DurationVarP(&scaleToZeroPodRetentionPeriod, "scale-to-zero-pod-retention-period", "", 0*time.Second, "the minimum seconds of time that the last pod will remain active after the Autoscaler has decided to scale pods to zero")
	AutoscalingUpdateCommand.Flags().StringVarP(&targetBurstCapacity, "target-burst-capacity", "", "200", "the desired burst capacity for the revision")
	AutoscalingUpdateCommand.Flags().StringVarP(&podAutoscalerClass, "pod-autoscaler-class", "", "kpa.autoscaling.knative.dev", "the config of Knative autoscaling to work with either the default KPA or a CPU based metric, i.e. Horizontal Pod Autoscaler (HPA)")
	AutoscalingUpdateCommand.Flags().StringVarP(&activatorCapacity, "activator-capacity", "", "200", "number of the concurrent requests an activator task can accept")

	return AutoscalingUpdateCommand
}
