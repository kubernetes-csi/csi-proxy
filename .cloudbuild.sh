#! /bin/bash

: ${CSI_PROW_BUILD_PLATFORMS:="windows amd64 .exe nanoserver:1809; windows amd64 .exe nanoserver:ltsc2022"}
: ${REGISTRY_NAME:="gcr.io/k8s-staging-sig-storage"}

. release-tools/prow.sh

gcr_cloud_build
