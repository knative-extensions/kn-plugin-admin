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
	knative.dev/client v0.29.1-0.20220204171521-6690a20e8f56
	knative.dev/hack v0.0.0-20220203160821-9b303d690fc9
	knative.dev/networking v0.0.0-20220208184716-c9e708231c88
	knative.dev/serving v0.29.1-0.20220208184517-52f07b02254b
)
