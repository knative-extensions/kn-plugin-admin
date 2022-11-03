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

import (
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"

	hprinters "knative.dev/client/pkg/printers"
)

// DomainListHandlers adds print handlers for domain list command
func DomainListHandlers(h hprinters.PrintHandler) {
	kDomainColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Custom-Domain", Type: "string", Description: "Name of Knative custom domain.", Priority: 1},
		{Name: "Selector", Type: "string", Description: "Selector of Knative custom domains.", Priority: 1},
	}
	h.TableHandler(kDomainColumnDefinitions, printKDomainList)
}

// printKDomainList populates the Knative custom domain list table rows
func printKDomainList(domainCM *corev1.ConfigMap, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	kDomainList := domainCM.Data
	delete(kDomainList, "_example")

	//sort the map for output
	sortedKeys := make([]string, 0, len(kDomainList))
	for k := range kDomainList {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	rows := make([]metav1beta1.TableRow, 0, len(kDomainList))
	for _, k := range sortedKeys {
		row := metav1beta1.TableRow{}
		row.Cells = append(row.Cells, k, formatSelectorForPrint(kDomainList[k]))
		rows = append(rows, []metav1beta1.TableRow{row}...)
	}
	return rows, nil
}

// format change for selector from "selector:\n  key1: value1\n  key2: value2\n" to "key1=value1; key2=value2;
func formatSelectorForPrint(selector string) string {
	parts := strings.Split(strings.ReplaceAll(strings.TrimSpace(selector), ":", "="), "\n")
	selectorForPrint := ""
	for i, v := range parts {
		//parts is from split by \n, so the first item i=0 will be start with "selector="", if not, return ""
		if i == 0 && !strings.HasPrefix(v, "selector=") {
			return ""

		} else if i > 0 { //skip the first item i=0 "selector="
			if strings.Contains(v, "=") {
				selectorForPrint = strings.Join([]string{selectorForPrint, strings.ReplaceAll(v, " ", "")}, "")
				//no ; for last selector entry
				if i < len(parts)-1 {
					selectorForPrint = strings.Join([]string{selectorForPrint, "; "}, "")
				}
			}
		}
	}
	return selectorForPrint
}
