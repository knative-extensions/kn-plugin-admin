module knative.dev/kn-plugin-admin

go 1.14

require (
	github.com/maximilien/kn-source-pkg v0.5.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/client v0.18.4
	knative.dev/serving v0.19.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.12
