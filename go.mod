module knative.dev/kn-plugin-admin

go 1.15

require (
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/maximilien/kn-source-pkg v0.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/client v0.21.1-0.20210301230649-270e32240152
	knative.dev/hack v0.0.0-20210203173706-8368e1f6eacf
	knative.dev/serving v0.21.1-0.20210302085014-cd7778194cb4
)

replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
