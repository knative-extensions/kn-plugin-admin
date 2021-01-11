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

package domain

import "testing"

func Test_formatSelectorForPrint(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"normal case with one selector key value", "selector:\n  key1: value1\n", "key1=value1"},
		{"normal case with two selector key value", "selector:\n  key1: value1\n  key2: value2\n", "key1=value1; key2=value2"},
		{"invalid input no selector", "notselector:\n  key1= value1\n", ""},
		{"invalid input no selector value", "selector:\n", ""},
		{"invalid input wrong selector value", "selector:\n  key1 value1\n", ""},
		{"invalid input wrong selector values", "selector:\n  key1 value1\n  key2: value2", "key2=value2"},
		{"empty selector", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := formatSelectorForPrint(tt.input)
			if output != tt.output {
				t.Errorf("formatSelectorForPrint() got = %v, want %v", output, tt.output)
			}
		})
	}
}
