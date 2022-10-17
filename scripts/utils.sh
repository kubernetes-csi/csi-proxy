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

compile_csi_proxy_integration_tests() {
  echo "Compiling CSI Proxy integration tests"
  GOOS=windows GOARCH=amd64 go test -c $pkgdir/integrationtests -o $pkgdir/bin/integrationtests.test.exe
}

sync_csi_proxy_integration_tests() {
  echo "Sync the integrationtests.exe binary"
  local integration_bin_path="$pkgdir/bin/integrationtests.test.exe"
  sync_file_to_vm $integration_bin_path
}

run_csi_proxy_integration_tests() {
  echo "Run integration tests"
  local ps1=$(cat << 'EOF'
  "& {
    $ErrorActionPreference = \"Stop\";
    .\integrationtests.test.exe --test.v
  }"
EOF
);

  gcloud compute ssh $windows_node --command="powershell -c $(echo $ps1 | tr '\n' ' ')"
}
