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
	knative.dev/client v0.26.1-0.20211101113922-db3bcbc6b04d
	knative.dev/hack v0.0.0-20211028194650-b96d65a5ff5e
	knative.dev/networking v0.0.0-20211029072251-c3606d9f7b38
	knative.dev/serving v0.26.1-0.20211101130423-44128f957caf
)
