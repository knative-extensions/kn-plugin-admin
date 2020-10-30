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

package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// LabelManagedBy is a label name to indicate who is managing this resource
var LabelManagedBy = "app.kubernetes.io/managed-by"

// AdminParams stores the configs for interacting with kube api
type AdminParams struct {
	KubeCfgPath        string
	ClientConfig       clientcmd.ClientConfig
	ClientSet          kubernetes.Interface
	InstallationMethod InstallationMethod
}

// InstallationMethod identify how knative get installed
type InstallationMethod int

const (
	// InstallationMethodUnknown default value
	InstallationMethodUnknown InstallationMethod = iota
	// InstallationMethodStandalone default installation method using full yaml configurations
	InstallationMethodStandalone
	// InstallationMethodOperator installation method using Knative Operator
	InstallationMethodOperator
)

// ErrorOperatorModeNotSupport indicates that knative is managed by operator and cannot handled by sub command
var ErrorOperatorModeNotSupport = errors.New("Knative managed by operator is not supported yet")

// ErrorInstallationMethodUnknown indicates that can not detect current installation method
var ErrorInstallationMethodUnknown = errors.New("Cannot detect current installation method")

// Initialize generate the clientset for params
func (params *AdminParams) Initialize() error {
	if params.ClientSet == nil {
		restConfig, err := params.RestConfig()
		if err != nil {
			return err
		}

		params.ClientSet, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			fmt.Println("failed to create client:", err)
			os.Exit(1)
		}
	}
	if params.InstallationMethod == InstallationMethodUnknown {
		im, err := params.installationMethod()
		if err != nil {
			fmt.Println("failed to get installation method:", err)
			os.Exit(1)
		}
		params.InstallationMethod = im
	}
	return nil

}

// installationMethod retrives the installation method
func (params *AdminParams) installationMethod() (InstallationMethod, error) {
	cm, err := params.ClientSet.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
	if err != nil {
		return InstallationMethodUnknown, err
	}
	for _, owner := range cm.OwnerReferences {
		if strings.HasPrefix(owner.APIVersion, "operator.knative.dev") && owner.Kind == "KnativeServing" {
			return InstallationMethodOperator, nil
		}
	}
	return InstallationMethodStandalone, nil
}

// EnsureInstallMethodStandalone return error if current installation method is not standalone
func (params *AdminParams) EnsureInstallMethodStandalone() error {
	switch params.InstallationMethod {
	case InstallationMethodOperator:
		return ErrorOperatorModeNotSupport
	case InstallationMethodUnknown:
		return ErrorInstallationMethodUnknown
	}
	return nil
}

// RestConfig returns REST config, which can be to use to create specific clientset
func (params *AdminParams) RestConfig() (*rest.Config, error) {
	var err error

	if params.ClientConfig == nil {
		params.ClientConfig, err = params.GetClientConfig()
		if err != nil {
			return nil, err
		}
	}

	config, err := params.ClientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetClientConfig gets ClientConfig from KubeCfgPath
func (params *AdminParams) GetClientConfig() (clientcmd.ClientConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(params.KubeCfgPath) == 0 {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}), nil
	}

	_, err := os.Stat(params.KubeCfgPath)
	if err == nil {
		loadingRules.ExplicitPath = params.KubeCfgPath
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}), nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	paths := filepath.SplitList(params.KubeCfgPath)
	if len(paths) > 1 {
		return nil, fmt.Errorf("Can not find config file. '%s' looks like a path. "+
			"Please use the env var KUBECONFIG if you want to check for multiple configuration files", params.KubeCfgPath)
	}
	return nil, fmt.Errorf("Config file '%s' can not be found", params.KubeCfgPath)
}
