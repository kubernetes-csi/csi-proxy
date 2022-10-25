# Development

## High Level Overview

CSI Proxy provides a suite of utility methods for running storage operations on Windows. Related functionalities are bundled into API groups packaged together, and storage drivers import individual API groups directly.

## Code structure

The CSI Proxy Go library is available under `pkg`.
  - `<API group>/<API group>.go` - The main entrypoint for an API group exposing an interface defining the API and a struct implementing that interface. Each API group takes in an underlying `HostAPI` implementation, which executes the actual implementation details, and the user facing interface only handles input/output validation.
  - `<API group>/types.go` - The request and response structs used by the API group's interface. Each method has a corresponding request and response type, defined in this file.
  - `<API group>/hostapi/hostapi.go` - The internal implementation of the CSI Proxy API for a particular API group.
  - `<API group>/hostapi/types.go` - Any types used by the particular `HostAPI` go here. Not all API groups have this file.
  - `utils/utils.go` - Shared utility functions between API groups.

There are unit tests scattered throughout `pkg`, but the main integration tests are defined `integrationtests`. Each API group has a corresponding test file, as well as a shared `utils.go` file.

## Making changes

Some notes about changing code:
- Version management is done using Go module. Please follow semantic versioning conventions when making API changes.
- Like all Kubernetes projects, CSI Proxy uses `go mod vendor`. Please run `go mod tidy && go mod vendor` to this code base's list of dependencies up to date.
- Update the unit tests in `pkg/<API group>/<API group>_test.go`.
- Update the integration tests in `integrationtests/`.

## Running E2E tests

There are a few presubmit tests that run in Github Actions, E2E tests need to run in a Windows VM with Hyper-V enabled, to run the E2E tests in Google Cloud follow this workflow:

- Create a Kubernetes e2e test cluster containing a Windows VM by following [these instructions](https://github.com/kubernetes/kubernetes/blob/master/cluster/gce/windows/README-GCE-Windows-kube-up.md) (follow step 2b when doing step 2). In addition to the flags given in the instructions, you also need to override/set:
```
export NUM_NODES=1
export NUM_WINDOWS_NODES=1
export NODE_LOCAL_SSDS_EXT="2,scsi,fs;2,nvme,fs"
export WINDOWS_ENABLE_HYPERV=true
export WINDOWS_NODE_OS_DISTRIBUTION=win2019
```
- Run the e2e tests using `scripts/run-integration.sh`.
