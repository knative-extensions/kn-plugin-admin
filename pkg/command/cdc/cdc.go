// Copyright 2021 The Knative Authors
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

package cdc

import (
	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-admin/pkg"
)

// NewCdcCommand return the CDC root command
func NewCdcCommand(p *pkg.AdminParams) *cobra.Command {
	var cdcCmd = &cobra.Command{
		Use:   "cdc",
		Short: "Manage cluster domain claim",
	}
	cdcCmd.AddCommand(NewCdcCreateCommand(p))
	cdcCmd.AddCommand(NewCdcListCommand(p))
	cdcCmd.AddCommand(NewCdcDeleteCommand(p))
	return cdcCmd
}
