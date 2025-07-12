module github.com/kubernetes-csi/csi-proxy

// NOTE: This project must be built with go < 1.23
// `make build` will error if go1.23 or higher is used.

go 1.24.3

godebug winsymlink=0

toolchain go1.24.4

require (
	github.com/Microsoft/go-winio v0.6.2
	github.com/go-ole/go-ole v1.3.0
	github.com/google/go-cmp v0.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.1
	github.com/iancoleman/strcase v0.3.0
	github.com/kubernetes-csi/csi-proxy/client v1.1.3
	github.com/microsoft/wmi v0.34.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.20.5
	github.com/sergi/go-diff v1.3.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.10.0
	golang.org/x/sys v0.32.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.0
	k8s.io/component-base v0.28.4
	k8s.io/gengo v0.0.0-20240911193312-2b36238f13e9
	k8s.io/klog/v2 v2.130.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241216192217-9240e9c98484 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.28.4 // indirect
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
