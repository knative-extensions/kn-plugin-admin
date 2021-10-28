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
	knative.dev/client v0.26.1-0.20211028082427-f027b38e200a
	knative.dev/hack v0.0.0-20211027200727-f1228dd5e3e6
	knative.dev/networking v0.0.0-20211027163221-b9fa70516e8e
	knative.dev/serving v0.26.1-0.20211028122826-e6b690a7f463
)
