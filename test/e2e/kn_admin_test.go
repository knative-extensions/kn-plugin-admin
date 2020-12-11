// Copyright 2020 The Knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	testcommon "github.com/maximilien/kn-source-pkg/test/e2e"
	"gotest.tools/assert"
	"knative.dev/client/lib/test"
	"knative.dev/client/pkg/util"
	"knative.dev/kn-plugin-admin/pkg/command/autoscaling"
	"knative.dev/serving/pkg/autoscaler/config"
)

const pluginName string = "admin"

type e2eTest struct {
	it         *testcommon.E2ETest
	kn         *test.Kn
	kubectl    *test.Kubectl
	backupData map[string]string
}

func newE2ETest(t *testing.T) *e2eTest {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil
	}

	it, err := testcommon.NewE2ETest("kn-admin", filepath.Join(currentDir, "../.."), false)
	if err != nil {
		return nil
	}

	kn := test.NewKn()
	kubectl := test.NewKubectl("knative-serving")
	e2eTest := &e2eTest{
		it:         it,
		kn:         &kn,
		kubectl:    &kubectl,
		backupData: make(map[string]string),
	}
	return e2eTest
}

func TestKnAdminPlugin(t *testing.T) {
	t.Parallel()

	e2eTest := newE2ETest(t)
	assert.Assert(t, e2eTest != nil)
	defer func() {
		assert.NilError(t, e2eTest.it.KnTest().Teardown())
	}()

	r := test.NewKnRunResultCollector(t, e2eTest.it.KnTest())
	defer r.DumpIfFailed()

	err := e2eTest.it.KnPlugin().Install()
	assert.NilError(t, err)

	t.Log("test kn admin domain subcommand")
	err = e2eTest.backupConfigMap("config-domain")
	assert.NilError(t, err)
	e2eTest.knAdminDomain(t, r)
	err = e2eTest.restoreConfigMap("config-domain")
	assert.NilError(t, err)

	t.Log("test kn admin registry subcommand")
	e2eTest.knAdminRegistry(t, r)
	err = e2eTest.backupConfigMap("config-observability")
	assert.NilError(t, err)

	t.Log("test kn admin autoscaling subcommand")
	err = e2eTest.backupConfigMap("config-autoscaler")
	assert.NilError(t, err)
	e2eTest.knAdminAutoscaling(t, r)
	err = e2eTest.restoreConfigMap("config-autoscaler")
	assert.NilError(t, err)

	t.Log("test kn admin profiling subcommand")
	e2eTest.knAdminProfiling(t, r)
	err = e2eTest.restoreConfigMap("config-observability")
	assert.NilError(t, err)

	err = e2eTest.it.KnPlugin().Uninstall()
	assert.NilError(t, err)
}

func (et *e2eTest) backupConfigMap(cm string) error {
	data, err := et.kubectl.Run("get", "configmap", cm, "-oyaml")
	if err != nil {
		return err
	}
	et.backupData[cm] = data
	return nil
}

func (et *e2eTest) restoreConfigMap(cm string) error {
	var data string
	var ok bool
	if data, ok = et.backupData[cm]; !ok {
		return fmt.Errorf("backup for configmap %s does not exists", cm)
	}
	f, err := ioutil.TempFile("", "")
	defer os.Remove(f.Name())
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	_, err = et.kubectl.Run("replace", "-f", f.Name(), "--force")
	return err
}

func (et *e2eTest) knAdminDomain(t *testing.T, r *test.KnRunResultCollector) {
	domain := "test.domain.test"
	out := et.kn.Run(pluginName, "domain", "set", "--custom-domain", domain)
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "domain", "set", "--custom-domain", domain, "--selector", "app=v1")
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "domain", "unset", "--custom-domain", domain)
	r.AssertNoError(out)
}

func (et *e2eTest) knAdminAutoscaling(t *testing.T, r *test.KnRunResultCollector) {
	out := et.kn.Run(pluginName, "autoscaling", "update", "--max-scale-up-rate", "2.5")
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "autoscaling", "update", "--stable-window", "2m", "--activator-capacity", "300")
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "autoscaling", "update", "--scale-to-zero")
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "autoscaling", "list")
	r.AssertNoError(out)

	// build autoscaler config data
	data := map[string]string{
		"max-scale-up-rate":    "2.5",
		"stable-window":        "2m",
		"activator-capacity":   "300",
		"enable-scale-to-zero": "true",
	}
	cfg, err := config.NewConfigFromMap(data)
	assert.NilError(t, err)

	// check config values in output
	lines := strings.Split(out.Stdout, "\n")
	for _, line := range lines {
		items := strings.Split(line, " ")
		if len(items) == 2 && items[0] != "NAME" {
			value := autoscaling.ConfigNameValueOfMap[items[0]](cfg)
			assert.Equal(t, value, items[1])
		}
	}
}

func (et *e2eTest) knAdminRegistry(t *testing.T, r *test.KnRunResultCollector) {
	out := et.kn.Run(pluginName, "registry", "add", "--username", "custom-user", "--password", "test", "--server", "test.test.io")
	r.AssertNoError(out)
	out = et.kn.Run(pluginName, "registry", "list")
	r.AssertNoError(out)
	outRows := strings.Split(out.Stdout, "\n")
	assert.Check(t, util.ContainsAll(outRows[0], "NAMESPACE", "SERVICEACCOUNT", "SECRET", "USERNAME", "SERVER", "EMAIL"))
	assert.Check(t, util.ContainsAll(outRows[1], "custom-user", "test.test.io"))
	out = et.kn.Run(pluginName, "registry", "remove", "--username", "custom-user", "--server", "test.test.io")
	r.AssertNoError(out)
}

func (et *e2eTest) knAdminProfiling(t *testing.T, r *test.KnRunResultCollector) {
	out := et.kn.Run(pluginName, "profiling", "--enable")
	r.AssertNoError(out)
	assert.Equal(t, "Knative Serving profiling is enabled\n", out.Stdout)

	out = et.kn.Run(pluginName, "profiling", "--heap", "--target", "controller")
	r.AssertNoError(out)
	assert.Check(t, strings.Contains(out.Stdout, "Saving heap profiling data to"))

	out = et.kn.Run(pluginName, "profiling", "--cpu", "5s", "--target", "controller")
	r.AssertNoError(out)
	assert.Check(t, strings.Contains(out.Stdout, "Saving 5 second(s) cpu profiling data to"))

	out = et.kn.Run(pluginName, "profiling", "--disable")
	r.AssertNoError(out)
	assert.Equal(t, "Knative Serving profiling is disabled\n", out.Stdout)
}
