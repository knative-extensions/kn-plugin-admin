module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	knative.dev/client v0.30.2-0.20220323165518-06c3202925e6
	knative.dev/hack v0.0.0-20220318020218-14f832e506f8
	knative.dev/networking v0.0.0-20220323170318-55757e9c20d6
	knative.dev/serving v0.30.1-0.20220325001548-ccc40402ae03
)
