module github.com/kubernetes-csi/csi-proxy

go 1.26.0

require (
	github.com/Microsoft/go-winio v0.6.2
	github.com/go-ole/go-ole v1.3.0
	github.com/google/go-cmp v0.7.0
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.1.0
	github.com/iancoleman/strcase v0.3.0
	github.com/kubernetes-csi/csi-proxy/client v1.1.3
	github.com/prometheus/client_golang v1.24.1
	github.com/sergi/go-diff v1.4.0
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.11.1
	golang.org/x/sys v0.47.0
	google.golang.org/grpc v1.83.0
	google.golang.org/protobuf v1.36.12-0.20260120151049-f2248ac996af
	k8s.io/component-base v0.36.3
	k8s.io/gengo v0.0.0-20240911193312-2b36238f13e9
	k8s.io/klog/v2 v2.140.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.36.3 // indirect
	k8s.io/klog v1.0.0 // indirect
)

replace (
	github.com/kubernetes-csi/csi-proxy/client => ./client

	// using my fork of gengo until
	// https://github.com/kubernetes/gengo/pull/155#issuecomment-537589085
	// is implemented, and the generic conversion generator merged into code-generator
	// FIXME: switch back to the upstream repo and/or code-generator!
	// (mauriciopoppe) while working on #140 I found out that I had to do an
	// override to the fork to stop generating auto* functions
	// https://github.com/mauriciopoppe/gengo/commit/9c78f58f3486e3c0cdb02ed9551d32762ac99773
	k8s.io/gengo => github.com/mauriciopoppe/gengo v0.0.0-20210525224835-9c78f58f3486
)
