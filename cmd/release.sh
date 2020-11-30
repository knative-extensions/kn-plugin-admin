#!/usr/bin/env bash
# Copyright © 2020 The Knative Authors
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

set -e
GIT_REVISION=`git rev-parse --short HEAD`
VERSION=`date -u '+v%Y%m%d'-`${GIT_REVISION}
BUILD_DATE=`date -u '+%Y-%m-%d %H:%M:%S'`
PKG="knative.dev/kn-plugin-admin/pkg/command"
LD_FLAGS="-X '${PKG}.Version=${VERSION}' -X '${PKG}.BuildDate=${BUILD_DATE}' -X '${PKG}.GitRevision=${GIT_REVISION}'"

export GO111MODULE=on
export CGO_ENABLED=0

# build for macOS
echo "🚧 🍏 Building for macOS"
GOOS=darwin GOARCH=amd64 go build -mod=readonly -ldflags "${LD_FLAGS}" -o kn-admin-darwin-amd64
# build for linux
echo "🚧 🐧 Building for Linux"
GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags "${LD_FLAGS}" -o kn-admin-linux-amd64
# build for windows
echo "🚧 🎠 Building for Windows"
GOOS=windows GOARCH=amd64 go build -mod=readonly -ldflags "${LD_FLAGS}" -o kn-admin-windows-amd64

