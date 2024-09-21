#!/bin/bash

# Importing this library shouldn't have side effects

pkgdir=${GOPATH}/src/github.com/kubernetes-csi/csi-proxy

# current_account is the current user
# in CI, the value is `prow`
current_account=$USER
# windows_node is the id of the GCE Windows instance
windows_node=$(kubectl get nodes -l kubernetes.io/os=windows -o jsonpath='{.items[*].metadata.name}')

# set the default zone for the gcloud sdk
if ! [ -z "${GCP_ZONE:-}" ]; then
  export CLOUDSDK_COMPUTE_ZONE="$GCP_ZONE"
fi

sync_file_to_vm() {
  gcloud compute scp $@ $windows_node:"C:\\Users\\${current_account}"
}

compile_csi_proxy() {
  echo "Compiling CSI Proxy"
  make -C $pkgdir build
}

compile_csi_proxy_integration_tests() {
  echo "Compiling CSI Proxy integration tests"
  GOOS=windows GOARCH=amd64 go test -c $pkgdir/integrationtests -o $pkgdir/bin/integrationtests.test.exe
}

sync_csi_proxy() {
  echo "Sync the csi-proxy.exe binary"
  local csi_proxy_bin_path="$pkgdir/bin/csi-proxy.exe"
  sync_file_to_vm $csi_proxy_bin_path
}

sync_csi_proxy_integration_tests() {
  echo "Sync the integrationtests.exe binary"
  local integration_bin_path="$pkgdir/bin/integrationtests.test.exe"
  sync_file_to_vm $integration_bin_path
}

sync_powershell_utils() {
  local utils_psm1="$pkgdir/scripts/utils.psm1"
  sync_file_to_vm $utils_psm1
}

restart_csi_proxy() {
  echo "Restart csi-proxy service"
  gcloud compute ssh $windows_node --command='powershell -c "& { $ErrorActionPreference = \"Stop\"; Import-Module (Resolve-Path(\"utils.psm1\")); Restart-CSIProxy; }"'
}

run_csi_proxy_integration_tests() {
  echo "Run integration tests"
  local ps1=$(cat << 'EOF'
  "& {
    $ErrorActionPreference = \"Stop\";
    Import-Module (Resolve-Path(\"utils.psm1\"));
    Run-CSIProxyIntegrationTests -test_args \"--test.v --test.run TestAPIGroups\";
    Run-CSIProxyIntegrationTests -test_args \"--test.v --test.run TestFilesystemAPIGroup\";
    Run-CSIProxyIntegrationTests -test_args \"--test.v --test.run TestDiskAPIGroup\";
    Run-CSIProxyIntegrationTests -test_args \"--test.v --test.run TestVolumeAPIs\";
    # Todo: Enable this test once the issue is fixed
    # Run-CSIProxyIntegrationTests -test_args \"--test.v --test.run TestSmbAPIGroup\";
  }"
EOF
);

  gcloud compute ssh $windows_node --command="powershell -c $(echo $ps1 | tr '\n' ' ')"
}
