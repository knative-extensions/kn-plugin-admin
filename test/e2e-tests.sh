#!/usr/bin/env bash

# Copyright 2020 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

source $(dirname $0)/common.sh

# Add local dir to have access to built kn
export PATH=$PATH:${REPO_ROOT_DIR}

# Will create and delete this namespace (used for all tests, modify if you want a different one)
export KN_E2E_NAMESPACE=kne2etests

export KNATIVE_EVENTING_VERSION="0.19.2"
export KNATIVE_SERVING_VERSION="0.19.0"

run() {
  # Create cluster
  initialize $@

  # Integration tests
  eval integration_test || fail_test

  success
}

integration_test() {
  header "Running kn-plugin-admin e2e tests for Knative Serving $KNATIVE_SERVING_VERSION and Eventing $KNATIVE_EVENTING_VERSION"
  go_test_e2e -timeout=45m ./test/e2e || return 1
}

# Fire up
run $@
