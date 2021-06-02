#/bin/bash

set -o nounset
set -ex

# the Google Storage bucket
: "${CSI_PROXY_BUCKET:?CSI_PROXY_BUCKET not set}"

# The bucket url of this script in Google Cloud, set in sync_scripts
SCRIPT_URL=

function sync_scripts {
  # upload initialization code to a bucket
  gsutil mb gs://${CSI_PROXY_BUCKET} || true
  SCRIPT_URL=gs://${CSI_PROXY_BUCKET}/e2e-runner.ps1
  gsutil cp scripts/e2e-runner.ps1 $SCRIPT_URL
}

function bootstrap_environment {
  if ! gcloud compute routes describe csi-proxy-windows-activation-route; then
    gcloud compute routes create csi-proxy-windows-activation-route \
      --destination-range=35.190.247.13/32 \
      --network=default \
      --next-hop-gateway=default-internet-gateway
  fi

  if ! gcloud compute firewall-rules describe csi-proxy-windows-activation-firewall-rule; then
    gcloud compute firewall-rules create csi-proxy-windows-activation-firewall-rule \
      --direction=EGRESS \
      --network=default \
      --action=ALLOW \
      --rules=tcp:1688 \
      --destination-ranges=35.190.247.13/32 \
      --priority=0
  fi

  if ! gcloud compute instances describe csi-proxy-e2e-tests; then
    gcloud compute instances create csi-proxy-e2e-tests \
      --image-project windows-cloud \
      --image-family windows-2019-core \
      --machine-type e2-medium \
      --boot-disk-size 100 \
      --boot-disk-type pd-ssd \
      --metadata=windows-startup-script-url=${SCRIPT_URL}
  fi
}

function main {
  sync_scripts
  bootstrap_environment
}

main
