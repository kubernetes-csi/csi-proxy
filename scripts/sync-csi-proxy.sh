#!/bin/bash

#
# Installs CSI Proxy in a kubernetes node
#
# Requirements:
# - a kubernetes cluster with a Windows nodepool
#
# Steps:
# - cross compile the binary
# - copy to the VM using scp
# - restart the CSI Proxy binary process with a helper powershell script

set -euxo pipefail

pkgdir=${GOPATH}/src/github.com/kubernetes-csi/csi-proxy
source $pkgdir/scripts/utils.sh

main() {
  compile_csi_proxy
  sync_csi_proxy
  sync_powershell_utils
  restart_csi_proxy
}

main
