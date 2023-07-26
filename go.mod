module github.com/kubernetes-csi/csi-proxy

go 1.20

require (
	github.com/Microsoft/go-winio v0.4.16
	github.com/google/go-cmp v0.5.0
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/kubernetes-csi/csi-proxy/client v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	golang.org/x/sys v0.10.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.25.0
	k8s.io/gengo v0.0.0-00010101000000-000000000000
	k8s.io/klog/v2 v2.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
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
