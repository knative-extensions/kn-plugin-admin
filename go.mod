module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/client v0.30.2-0.20220314132918-d75836c4f42a
	knative.dev/hack v0.0.0-20220314052818-c9c3ea17a2e9
	knative.dev/networking v0.0.0-20220316064759-31d0d3ffe54c
	knative.dev/serving v0.30.1-0.20220315154224-8fe13873eb54
)
