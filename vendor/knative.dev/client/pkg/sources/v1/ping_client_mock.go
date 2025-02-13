// Copyright © 2019 The Knative Authors
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

package v1

import (
	"context"
	"testing"

	"knative.dev/client/pkg/util/mock"
	sourcesv1 "knative.dev/eventing/pkg/apis/sources/v1"
)

type MockKnPingSourceClient struct {
	t        *testing.T
	recorder *PingSourcesRecorder
}

// NewMockKnPingSourceClient returns a new mock instance which you need to record for
func NewMockKnPingSourceClient(t *testing.T, ns ...string) *MockKnPingSourceClient {
	namespace := "default"
	if len(ns) > 0 {
		namespace = ns[0]
	}
	return &MockKnPingSourceClient{
		t:        t,
		recorder: &PingSourcesRecorder{mock.NewRecorder(t, namespace)},
	}
}

// Ensure that the interface is implemented
var _ KnPingSourcesClient = &MockKnPingSourceClient{}

// recorder for service
type PingSourcesRecorder struct {
	r *mock.Recorder
}

// Recorder returns the recorder for registering API calls
func (c *MockKnPingSourceClient) Recorder() *PingSourcesRecorder {
	return c.recorder
}

// Namespace of this client
func (c *MockKnPingSourceClient) Namespace() string {
	return c.recorder.r.Namespace()
}

// CreatePingSource records a call for CreatePingSource with the expected error
func (sr *PingSourcesRecorder) CreatePingSource(pingSource interface{}, err error) {
	sr.r.Add("CreatePingSource", []interface{}{pingSource}, []interface{}{err})
}

// CreatePingSource performs a previously recorded action, failing if non has been registered
func (c *MockKnPingSourceClient) CreatePingSource(ctx context.Context, pingSource *sourcesv1.PingSource) error {
	call := c.recorder.r.VerifyCall("CreatePingSource", pingSource)
	return mock.ErrorOrNil(call.Result[0])
}

// GetPingSource records a call for GetPingSource with the expected object or error. Either pingsource or err should be nil
func (sr *PingSourcesRecorder) GetPingSource(name interface{}, pingSource *sourcesv1.PingSource, err error) {
	sr.r.Add("GetPingSource", []interface{}{name}, []interface{}{pingSource, err})
}

// GetPingSource performs a previously recorded action, failing if non has been registered
func (c *MockKnPingSourceClient) GetPingSource(ctx context.Context, name string) (*sourcesv1.PingSource, error) {
	call := c.recorder.r.VerifyCall("GetPingSource", name)
	return call.Result[0].(*sourcesv1.PingSource), mock.ErrorOrNil(call.Result[1])
}

// UpdatePingSource records a call for UpdatePingSource with the expected error (nil if none)
func (sr *PingSourcesRecorder) UpdatePingSource(pingSource interface{}, err error) {
	sr.r.Add("UpdatePingSource", []interface{}{pingSource}, []interface{}{err})
}

// UpdatePingSource performs a previously recorded action, failing if non has been registered
func (c *MockKnPingSourceClient) UpdatePingSource(ctx context.Context, pingSource *sourcesv1.PingSource) error {
	call := c.recorder.r.VerifyCall("UpdatePingSource", pingSource)
	return mock.ErrorOrNil(call.Result[0])
}

func (c *MockKnPingSourceClient) UpdatePingSourceWithRetry(ctx context.Context, name string, updateFunc PingSourceUpdateFunc, nrRetries int) error {
	return updatePingSourceWithRetry(ctx, c, name, updateFunc, nrRetries)
}

// UpdatePingSource records a call for DeletePingSource with the expected error (nil if none)
func (sr *PingSourcesRecorder) DeletePingSource(name interface{}, err error) {
	sr.r.Add("DeletePingSource", []interface{}{name}, []interface{}{err})
}

// DeletePingSource performs a previously recorded action, failing if non has been registered
func (c *MockKnPingSourceClient) DeletePingSource(ctx context.Context, name string) error {
	call := c.recorder.r.VerifyCall("DeletePingSource", name)
	return mock.ErrorOrNil(call.Result[0])
}

// ListPingSource records a call for ListPingSource with the expected error (nil if none)
func (sr *PingSourcesRecorder) ListPingSource(pingSourceList *sourcesv1.PingSourceList, err error) {
	sr.r.Add("ListPingSource", []interface{}{}, []interface{}{pingSourceList, err})
}

// ListPingSource performs a previously recorded action, failing if non has been registered
func (c *MockKnPingSourceClient) ListPingSource(context.Context) (*sourcesv1.PingSourceList, error) {
	call := c.recorder.r.VerifyCall("ListPingSource")
	return call.Result[0].(*sourcesv1.PingSourceList), mock.ErrorOrNil(call.Result[1])
}

// Validates validates whether every recorded action has been called
func (sr *PingSourcesRecorder) Validate() {
	sr.r.CheckThatAllRecordedMethodsHaveBeenCalled()
}
