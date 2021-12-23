#/bin/bash
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

set -ex

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
csi_proxy_bin_path="$script_dir/../bin/csi-proxy.exe"
sync_script_ps_path="$script_dir/sync-csi-proxy.ps1"

main() {
  echo "Compiling CSI Proxy"
  make build

  local windows_node=$(kubectl get nodes -l kubernetes.io/os=windows -o jsonpath='{.items[*].metadata.name}')
  gcloud compute instances describe $windows_node > /dev/null

  echo "Sync the csi-proxy.exe binary"
  local current_account=$(gcloud config list account --format "value(core.account)" | sed -r 's/@\S+//g')
  gcloud compute scp $csi_proxy_bin_path $sync_script_ps_path $windows_node:"C:\\Users\\${current_account}"

  echo "Restart csi-proxy service"
  gcloud compute ssh $windows_node --command="powershell .\sync-csi-proxy.ps1"
}

main
