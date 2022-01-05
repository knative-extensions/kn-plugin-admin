module knative.dev/kn-plugin-admin

go 1.16

require (
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/client v0.28.1-0.20220104123133-63142983acd3
	knative.dev/hack v0.0.0-20211222071919-abd085fc43de
	knative.dev/networking v0.0.0-20211223134928-e40187c3026d
	knative.dev/serving v0.28.1-0.20220104122631-278af32f24ce
)
