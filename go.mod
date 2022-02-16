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
	knative.dev/client v0.29.1-0.20220215151759-24b6184341b7
	knative.dev/hack v0.0.0-20220215185059-b9cb1983b600
	knative.dev/networking v0.0.0-20220216014839-4337f034f4ca
	knative.dev/serving v0.29.1-0.20220216014040-0ee7b6f0b49e
)
