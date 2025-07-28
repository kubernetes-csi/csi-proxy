module github.com/kubernetes-csi/csi-proxy/v2

go 1.24.3

godebug winsymlink=0

toolchain go1.24.4

require (
	github.com/go-ole/go-ole v1.3.0
	github.com/microsoft/wmi v0.34.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.32.0
	k8s.io/klog/v2 v2.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)
