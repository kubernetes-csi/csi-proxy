#!/bin/bash

# A Prow job can override these defaults, but this shouldn't be necessary.

# # Only these tests make sense for csi-proxy
: ${CSI_PROW_TESTS:="unit"}
: ${CSI_PROW_BUILD_PLATFORMS:="windows amd64 .exe"}

 . release-tools/prow.sh

main
