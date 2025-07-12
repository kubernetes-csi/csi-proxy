#!/bin/bash

#
# Runs the integration tests
#
# Requirements:
# - a kubernetes cluster with a Windows nodepool
#
# Steps:
# - cross compile the csi-proxy binary and the integration tests
# - copy to the VM using scp
# - restart the CSI Proxy binary process with a helper powershell script
# - run the integration tests

set -euxo pipefail

pkgdir=${GOPATH}/src/github.com/kubernetes-csi/csi-proxy
source $pkgdir/scripts/utils.sh

main() {
  compile_csi_proxy
  compile_csi_proxy_integration_tests
  sync_csi_proxy
  sync_csi_proxy_integration_tests
  sync_powershell_utils
  restart_csi_proxy
  run_csi_proxy_integration_tests
}

main
