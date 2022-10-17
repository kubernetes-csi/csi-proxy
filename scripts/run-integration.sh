#!/bin/bash

#
# Runs the integration tests
#
# Requirements:
# - a kubernetes cluster with a Windows nodepool
#
# Steps:
# - cross compile the csi-proxy integration tests
# - copy to the VM using scp
# - run the integration tests

set -euxo pipefail

pkgdir=${GOPATH}/src/github.com/kubernetes-csi/csi-proxy
source $pkgdir/scripts/utils.sh

main() {
  compile_csi_proxy_integration_tests
  sync_csi_proxy_integration_tests
  run_csi_proxy_integration_tests
}

main
