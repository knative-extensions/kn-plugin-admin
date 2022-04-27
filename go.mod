module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.1.0
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/client v0.31.1-0.20220427131752-2eaeab491e38
	knative.dev/hack v0.0.0-20220427014036-5f473869d377
	knative.dev/networking v0.0.0-20220427030951-700c5762af5b
	knative.dev/serving v0.31.1-0.20220427014148-6b3905bfa89f
)
