#!/bin/bash

# A Prow job can override these defaults, but this shouldn't be necessary.

# # Only these tests make sense for csi-proxy
: ${CSI_PROW_TESTS:="unit"}
: ${CSI_PROW_BUILD_PLATFORMS:="windows amd64 .exe nanoserver:1809"}

. release-tools/prow.sh

# This creates the CSI_PROW_WORK directory that is needed by run_with_go.
ensure_paths

# main
run_with_go "${CSI_PROW_GO_VERSION_BUILD}" make all "GOFLAGS_VENDOR=${GOFLAGS_VENDOR}" "BUILD_PLATFORMS=${CSI_PROW_BUILD_PLATFORMS}"
run_with_go "${CSI_PROW_GO_VERSION_BUILD}" make -k test "GOFLAGS_VENDOR=${GOFLAGS_VENDOR}" 2>&1 | make_test_to_junit

# build / push multi-arch images for validation
gcr_cloud_build