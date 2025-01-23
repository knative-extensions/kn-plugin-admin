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

package e2e

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"knative.dev/client/pkg/util/test"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type E2ETest struct {
	knTest   *test.KnTest
	knPlugin *knPlugin
}

// NewE2ETest for pluginName in pluginPath
func NewE2ETest(pluginName string, pluginPath string, install bool) (*E2ETest, error) {
	knTest, err := test.NewKnTest()
	if err != nil {
		return nil, err
	}

	knPlugin := &knPlugin{
		kn:         knTest.Kn(),
		pluginName: pluginName,
		pluginPath: pluginPath,
		install:    install,
	}

	e2eTest := &E2ETest{
		knTest:   knTest,
		knPlugin: knPlugin,
	}

	return e2eTest, nil
}

// KnTest object
func (e2eTest *E2ETest) KnTest() *test.KnTest {
	return e2eTest.knTest
}

// KnPlugin object
func (e2eTest *E2ETest) KnPlugin() *knPlugin {
	return e2eTest.knPlugin
}

type knPlugin struct {
	kn         test.Kn
	pluginName string
	pluginPath string
	install    bool
}

// Run the KnPlugin returning a KnRunResult
func (kp *knPlugin) Run(args ...string) test.KnRunResult {
	if kp.install {
		err := kp.Install()
		if err != nil {
			fmt.Printf("error installing kn plugin: %s\n", err.Error())
			return test.KnRunResult{}
		}
		defer kp.Uninstall()
	}
	return RunKnPlugin(kp.kn.Namespace(), kp.pluginName, args)
}

// Kn object to run `kn`
func (kp *knPlugin) Kn() test.Kn {
	return kp.kn
}

// Install the KnPlugin
func (kp *knPlugin) Install() error {
	configDir, err := defaultConfigDir()
	if err != nil {
		fmt.Printf("error determining config directory: %s\n", err.Error())
		return err
	}

	pluginDir := filepath.Join(configDir, "plugins")
	if !dirExists(pluginDir) {
		err = os.MkdirAll(pluginDir, 0700)
		if err != nil {
			return err
		}
	}

	fmt.Printf("installing 'kn' plugin '%s' from path '%s' to config path: %s\n", kp.pluginName, kp.pluginPath, configDir)

	err = copyPluginFile(filepath.Join(kp.pluginPath, kp.pluginName), filepath.Join(pluginDir, kp.pluginName))
	if err != nil {
		fmt.Printf("error copying plugin file to config directory: %s\n", err.Error())
		return err
	}

	return nil
}

// Uninstall the KnPlugin
func (kp *knPlugin) Uninstall() error {
	configDir, err := defaultConfigDir()
	if err != nil {
		fmt.Printf("error determining config directory: %s\n", err.Error())
		return err
	}

	fmt.Printf("uninstalling 'kn' plugin '%s' from config path '%s'\n", kp.pluginName, configDir)

	err = os.Remove(filepath.Join(configDir, "plugins", kp.pluginName))
	if err != nil {
		fmt.Printf("error removing plugin from config directory: %s\n", err.Error())
		return err
	}

	return nil
}

// Utility functions

func copyPluginFile(sourceFile string, destDir string) error {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(destDir, input, 0700)
	if err != nil {
		return err
	}

	return nil
}

func defaultConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// Check the deprecated path first and fallback to it, add warning to error message
	if configHome := filepath.Join(home, ".kn"); dirExists(configHome) {
		migrationPath := filepath.Join(home, ".config", "kn")
		if runtime.GOOS == "windows" {
			migrationPath = filepath.Join(os.Getenv("APPDATA"), "kn")
		}
		return configHome, fmt.Errorf("WARNING: deprecated kn config directory detected. "+
			"Please move your configuration to: %s", migrationPath)
	}
	// Respect %APPDATA% on MS Windows
	// C:\Documents and Settings\username\Application JsonData
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "kn"), nil
	}
	// Respect XDG_CONFIG_HOME if set
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return filepath.Join(xdgHome, "kn"), nil
	}
	// Fallback to XDG default for both Linux and macOS
	// ~/.config/kn
	return filepath.Join(home, ".config", "kn"), nil
}

func dirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func pluginArgs(pluginName string) []string {
	pluginParts := strings.Split(pluginName, "-")
	return pluginParts[1:]
}

func RunKnPlugin(namespace string, pluginName string, args []string) test.KnRunResult {
	pluginArgs := pluginArgs(pluginName)
	args = append(args, []string{"--namespace", namespace}...)
	argsWithPlugin := append(pluginArgs, args...)
	return test.RunKn(namespace, argsWithPlugin)
}
