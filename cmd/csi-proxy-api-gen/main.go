package main

import (
	"os"

	"github.com/kubernetes-csi/csi-proxy/cmd/csi-proxy-api-gen/generators"
)

func main() {
	generators.Execute(os.Args[0], os.Args[1:]...)
}
