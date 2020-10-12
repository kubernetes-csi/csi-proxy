module github.com/kubernetes-csi/csi-proxy

go 1.13

replace (
	github.com/kubernetes-csi/csi-proxy/client => ./client

	// using my fork of gengo until
	// https://github.com/kubernetes/gengo/pull/155#issuecomment-537589085
	// is implemented, and the generic conversion generator merged into code-generator
	// FIXME: switch back to the upstream repo and/or code-generator!
	k8s.io/gengo => github.com/wk8/gengo v0.0.0-20191007012548-3d2530bfe606
)

require (
	github.com/Microsoft/go-winio v0.4.14
	github.com/golang/protobuf v1.4.2
	github.com/iancoleman/strcase v0.1.2
	github.com/kubernetes-csi/csi-proxy/client v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634
	golang.org/x/tools v0.0.0-20201011145850-ed2f50202694 // indirect
	google.golang.org/genproto v0.0.0-20201009135657-4d944d34d83c // indirect
	google.golang.org/grpc v1.33.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	k8s.io/gengo v0.0.0-20200728071708-7794989d0000
	k8s.io/klog v1.0.0
)
