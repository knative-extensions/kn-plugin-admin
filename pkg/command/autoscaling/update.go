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

// Autoscaling global configs
type Config struct {
	ScaleToZero                          bool
	EnableScaleToZero                    string
	StableWindow                         time.Duration
	ScaleToZeroGracePeriod               time.Duration
	ScaleToZeroPodRetentionPeriod        time.Duration
	ContainerConcurrencyTargetPercentage string
	PanicWindowPercentage                string
	PanicThresholdPercentage             string
	MaxScaleUpRate                       string
	MaxScaleDownRate                     string
	TargetBurstCapacity                  string
	ActivatorCapacity                    string
	RequestsPerSecondTargetDefault       string
	ContainerConcurrencyTargetDefault    string
	PodAutoscalerClass                   string
}

func NewConfig() Config {
	config := Config{}
	config.EnableScaleToZero = "enable-scale-to-zero"
	return config
}

var (
	knativeServing   = "knative-serving"
	configAutoscaler = "config-autoscaler"
)

func NewAutoscalingUpdateCommand(p *pkg.AdminParams) *cobra.Command {
	config := NewConfig()
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
			client, err := p.NewKubeClient()
			if err != nil {
				return err
			}

			currentCm := &corev1.ConfigMap{}
			currentCm, err = client.CoreV1().ConfigMaps(knativeServing).Get(context.TODO(), configAutoscaler, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get ConfigMaps: %+v", err)
			}
			desiredCm := currentCm.DeepCopy()

			if cmd.Flags().Changed("scale-to-zero") && cmd.Flags().Changed("no-scale-to-zero") {
				return fmt.Errorf("please specify either --scale-to-zero or --no-scale-to-zero")
			}

			if cmd.Flags().Changed("scale-to-zero") {
				desiredCm.Data[config.EnableScaleToZero] = "true"
			}

			if cmd.Flags().Changed("no-scale-to-zero") {
				desiredCm.Data[config.EnableScaleToZero] = "false"
			}

			if cmd.Flags().Changed("requests-per-second-target-default") {
				desiredCm.Data["requests-per-second-target-default"] = fmt.Sprintf("%s", config.RequestsPerSecondTargetDefault)
			}

			if cmd.Flags().Changed("container-concurrency-target-default") {
				desiredCm.Data["container-concurrency-target-default"] = fmt.Sprintf("%s", config.ContainerConcurrencyTargetDefault)
			}

			if cmd.Flags().Changed("container-concurrency-target-percentage") {
				desiredCm.Data["container-concurrency-target-percentage"] = fmt.Sprintf("%s", config.ContainerConcurrencyTargetPercentage)
			}

			if cmd.Flags().Changed("stable-window") {
				if config.StableWindow < as.WindowMin || config.StableWindow > as.WindowMax {
					return fmt.Errorf("stable-window = %v, must be in [%v; %v] range", config.StableWindow,
						as.WindowMin, as.WindowMax)
				}

				if config.StableWindow.Round(time.Second) != config.StableWindow {
					return fmt.Errorf("stable-window = %v, must be specified with at most second precision", config.StableWindow)
				}

				fmt.Printf("debug stable-window %vs\n", config.StableWindow.Seconds())
				desiredCm.Data["stable-window"] = fmt.Sprintf("%vs", config.StableWindow.Seconds())
			}

			if cmd.Flags().Changed("panic-window-percentage") {
				desiredCm.Data["panic-window-percentage"] = config.PanicWindowPercentage
			}

			if cmd.Flags().Changed("panic-threshold-percentage") {
				desiredCm.Data["panic-threshold-percentage"] = config.PanicThresholdPercentage
			}

			if cmd.Flags().Changed("max-scale-up-rate") {
				tmp, err := strconv.ParseFloat(config.MaxScaleUpRate, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", config.MaxScaleUpRate)
				}
				if tmp <= 1.0 {
					return fmt.Errorf("max-scale-up-rate = %v, must be greater than 1.0", config.MaxScaleUpRate)
				}
				desiredCm.Data["max-scale-up-rate"] = config.MaxScaleUpRate
			}

			if cmd.Flags().Changed("max-scale-down-rate") {
				tmp, err := strconv.ParseFloat(config.MaxScaleDownRate, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", config.MaxScaleUpRate)
				}
				if tmp <= 1.0 {
					return fmt.Errorf("max-scale-down-rate = %v, must be greater than 1.0", config.MaxScaleDownRate)
				}
				desiredCm.Data["max-scale-down-rate"] = config.MaxScaleDownRate
			}

			if cmd.Flags().Changed("scale-to-zero-grace-period") {
				if config.ScaleToZeroGracePeriod < as.WindowMin {
					return fmt.Errorf("scale-to-zero-grace-period must be at least %v, got %v", as.WindowMin, config.ScaleToZeroGracePeriod)
				}

				desiredCm.Data["scale-to-zero-grace-period"] = fmt.Sprintf("%vs", config.ScaleToZeroGracePeriod.Seconds())
			}

			if cmd.Flags().Changed("scale-to-zero-pod-retention-period") {
				if config.ScaleToZeroPodRetentionPeriod < 0 {
					return fmt.Errorf("scale-to-zero-pod-retention-period cannot be negative, was: %v", config.ScaleToZeroPodRetentionPeriod)
				}
				desiredCm.Data["scale-to-zero-pod-retention-period"] = fmt.Sprintf("%vs", config.ScaleToZeroPodRetentionPeriod.Seconds())
			}

			if cmd.Flags().Changed("target-burst-capacity") {
				tmp, err := strconv.ParseFloat(config.TargetBurstCapacity, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", config.MaxScaleUpRate)
				}
				if tmp < 0 && tmp != -1 {
					return fmt.Errorf("target-burst-capacity must be either non-negative or -1 (for unlimited), got %s", config.TargetBurstCapacity)
				}
				desiredCm.Data["target-burst-capacity"] = config.TargetBurstCapacity
			}

			if cmd.Flags().Changed("pod-autoscaler-class") {
				desiredCm.Data["pod-autoscaler-class"] = config.PodAutoscalerClass
			}

			if cmd.Flags().Changed("activator-capacity") {
				tmp, err := strconv.ParseFloat(config.ActivatorCapacity, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %v", config.MaxScaleUpRate)
				}
				if tmp < 1 {
					return fmt.Errorf("activator-capacity = %v, must be at least 1", config.ActivatorCapacity)
				}
				desiredCm.Data["activator-capacity"] = config.ActivatorCapacity
			}

			err = utils.UpdateConfigMap(client, desiredCm)
			if err != nil {
				return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", configAutoscaler, knativeServing, err)
			}
			cmd.Printf("Updated Knative autoscaling config\n")

			return nil
		},
	}

	flags.AddBothBoolFlagsUnhidden(AutoscalingUpdateCommand.Flags(), &config.ScaleToZero, "scale-to-zero", "", true,
		"Enable scale-to-zero if set.")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.RequestsPerSecondTargetDefault, "requests-per-second-target-default", "", "200", "the default target value for requests per second")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.ContainerConcurrencyTargetDefault, "container-concurrency-target-default", "", "100", "the default value of container concurrency target")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.ContainerConcurrencyTargetPercentage, "container-concurrency-target-percentage", "", "0.7", "percentage of the specified target should actually be targeted by the Autoscaler")
	AutoscalingUpdateCommand.Flags().DurationVarP(&config.StableWindow, "stable-window", "", 60*time.Second, "when operating in a stable mode, the autoscaler operates on the average concurrency over the x seconds of stable window")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.PanicWindowPercentage, "panic-window-percentage", "", "10", "The panic window is defined as a percentage of the stable window")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.PanicThresholdPercentage, "panic-threshold-percentage", "", "200", "This threshold defines when the autoscaler will move from stable mode into panic mode")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.MaxScaleUpRate, "max-scale-up-rate", "", "1000", "Maximum ratio of desired vs. observed pods")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.MaxScaleDownRate, "max-scale-down-rate", "", "2", "Maximum ratio of observed vs. desired pods")
	AutoscalingUpdateCommand.Flags().DurationVarP(&config.ScaleToZeroGracePeriod, "scale-to-zero-grace-period", "", 30*time.Second, "the maximum seconds of time that the last pod will remain active after the Autoscaler has decided to scale pods to zero")
	AutoscalingUpdateCommand.Flags().DurationVarP(&config.ScaleToZeroPodRetentionPeriod, "scale-to-zero-pod-retention-period", "", 0*time.Second, "the minimum seconds of time that the last pod will remain active after the Autoscaler has decided to scale pods to zero")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.TargetBurstCapacity, "target-burst-capacity", "", "200", "the desired burst capacity for the revision")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.PodAutoscalerClass, "pod-autoscaler-class", "", "kpa.autoscaling.knative.dev", "the config of Knative autoscaling to work with either the default KPA or a CPU based metric, i.e. Horizontal Pod Autoscaler (HPA)")
	AutoscalingUpdateCommand.Flags().StringVarP(&config.ActivatorCapacity, "activator-capacity", "", "200", "number of the concurrent requests an activator task can accept")

	return AutoscalingUpdateCommand
}
