/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1alpha2 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha2"
)

type FakeSourcesV1alpha2 struct {
	*testing.Fake
}

func (c *FakeSourcesV1alpha2) ApiServerSources(namespace string) v1alpha2.ApiServerSourceInterface {
	return &FakeApiServerSources{c, namespace}
}

func (c *FakeSourcesV1alpha2) ContainerSources(namespace string) v1alpha2.ContainerSourceInterface {
	return &FakeContainerSources{c, namespace}
}

func (c *FakeSourcesV1alpha2) PingSources(namespace string) v1alpha2.PingSourceInterface {
	return &FakePingSources{c, namespace}
}

func (c *FakeSourcesV1alpha2) SinkBindings(namespace string) v1alpha2.SinkBindingInterface {
	return &FakeSinkBindings{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSourcesV1alpha2) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
