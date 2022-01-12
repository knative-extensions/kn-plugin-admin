module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/client v0.28.1-0.20220111130713-1fcab77a0876
	knative.dev/hack v0.0.0-20220111151514-59b0cf17578e
	knative.dev/networking v0.0.0-20220112013650-eac673fb5c49
	knative.dev/serving v0.28.1-0.20220112024650-39af71665e8f
)
