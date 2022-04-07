module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/client v0.30.2-0.20220405102843-73b2d8bb1370
	knative.dev/hack v0.0.0-20220401031746-a75ca495e7f4
	knative.dev/networking v0.0.0-20220407031944-7fa8012b6f2d
	knative.dev/serving v0.30.1-0.20220407131646-c2e849153a93
)
