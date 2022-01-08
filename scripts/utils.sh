#/bin/bash
# Importing this library shouldn't have side effects

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# current_account is the current Google SDK Account
current_account=$(gcloud config list account --format "value(core.account)" | sed -r 's/@\S+//g')
# windows_node is the id of the GCE Windows instance
windows_node=$(kubectl get nodes -l kubernetes.io/os=windows -o jsonpath='{.items[*].metadata.name}')

if ! [ -z $GCP_ZONE ]; then
  export CLOUDSDK_COMPUTE_ZONE=$GCP_ZONE
fi

sync_file_to_vm() {
  gcloud compute scp $@ $windows_node:"C:\\Users\\${current_account}"
}

compile_csi_proxy() {
  echo "Compiling CSI Proxy"
  make build
}

compile_csi_proxy_integration_tests() {
  echo "Compiling CSI Proxy integration tests"
  GOOS=windows GOARCH=amd64 go test -c ./integrationtests -o bin/integrationtests.test.exe
}

sync_csi_proxy() {
  echo "Sync the csi-proxy.exe binary"
  local csi_proxy_bin_path="$script_dir/../bin/csi-proxy.exe"
  sync_file_to_vm $csi_proxy_bin_path
}

sync_csi_proxy_integration_tests() {
  echo "Sync the integrationtests.exe binary"
  local integration_bin_path="$script_dir/../bin/integrationtests.test.exe"
  sync_file_to_vm $integration_bin_path
}

sync_powershell_utils() {
  local utils_psm1="$script_dir/utils.psm1"
  sync_file_to_vm $utils_psm1
}

restart_csi_proxy() {
  echo "Restart csi-proxy service"
  gcloud compute ssh $windows_node --command='powershell -c "& { Import-Module .\utils.psm1; Restart-CSIProxy }"'
}

run_csi_proxy_integration_tests() {
  echo "Run integration tests"
  gcloud compute ssh $windows_node --command='powershell -c "& { Import-Module .\utils.psm1; Run-CSIProxyIntegrationTests }"'
}
