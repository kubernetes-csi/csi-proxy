#! /bin/bash

. release-tools/prow.sh

# Extract tag-n-hash value from GIT_TAG (form vYYYYMMDD-tag-n-hash) for REV value.
REV=v$(echo $GIT_TAG | cut -f3- -d 'v')

run_with_go "${CSI_PROW_GO_VERSION_BUILD}" make build REV=${REV}
cp bin/csi-proxy.exe bin/csi-proxy-${PULL_BASE_REF}.exe
