#! /bin/bash

. release-tools/prow.sh

# Extract tag-n-hash value from GIT_TAG (form vYYYYMMDD-tag-n-hash) for REV value.
REV=v$(echo "$GIT_TAG" | cut -f3- -d 'v')

# This creates the CSI_PROW_WORK directory that is needed by run_with_go.
ensure_paths

run_with_go "${CSI_PROW_GO_VERSION_BUILD}" make build REV="${REV}"
