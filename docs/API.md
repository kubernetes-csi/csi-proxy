# CSI-proxy's API

## Overview

CSI-proxy's API is a GRPC, versioned API.

The server exposes a number of API groups, all independent of each other. Additionally, each API group has one or more versions. Each version in each group listens for GRPC messages on a Windows named pipe of the form `\\.\\pipe\\csi-proxy-<api_group_name>-<version>` (e.g. `\\.\\pipe\\csi-proxy-filesystem-v2alpha1` or `\\.\\pipe\\csi-proxy-iscsi-v1`).

APIs are defined by protobuf files; each API group should live in its own directory under `client/api/<api_group_name>` in this repo's root (e.g. `client/api/iscsi`), and then define each of its version in `client/api/<api_group_name>/<version>/api.proto` files (e.g. `client/api/iscsi/v1/api.proto`). Each `proto` file should define exactly one RPC service.

Internally, there is only one server `struct` per API group, that handles all the versions for that API group. That server is defined in this repo's `pkg/server/<api_group_name>` (e.g. `pkg/server/iscsi`) go package. This go package should follow the following pattern:

<a name="serverPkgTree"></a>
```
pkg/server/<api_group_name>
├── impl
│   └── types.go
└── server.go
```
where `types.go` should contain the internal types corresponding to the various protobuf types for that API group - these internal structures must be able to represent all the different versions of the API. For example, given a `dummy` API group with two versions defined by the following `proto` files:

`client/api/dummy/v1alpha1/api.proto`
```proto
syntax = "proto3";

package v1alpha1;

service Dummy {
    // ComputeDouble computes the double of the input. Real smart stuff!
    rpc ComputeDouble(ComputeDoubleRequest) returns (ComputeDoubleResponse) {}
}

message ComputeDoubleRequest{
    int32 input32 = 1;
}

message ComputeDoubleResponse{
    int32 response32 = 1;
}
```

and

`client/api/dummy/v1/api.proto`
```proto
syntax = "proto3";

package v1;

service Dummy {
    // ComputeDouble computes the double of the input. Real smart stuff!
    rpc ComputeDouble(ComputeDoubleRequest) returns (ComputeDoubleResponse) {}
}

message ComputeDoubleRequest{
    // we changed in favor of an int64 field here
    int64 input = 2;
}

message ComputeDoubleResponse{
    int64 response = 2;

    // set to true if the result overflowed
    bool overflow = 3;
}
```

then `pkg/server/dummy/impl/types.go` could look something like:
```go
type ComputeDoubleRequest struct {
	Input int64
}

type ComputeDoubleResponse struct {
	Response int64
	Overflow bool
}
```
and then the API group's server (`pkg/server/dummy/server.go`) needs to define the callbacks to handle requests for all API versions, e.g.:
```go
type Server struct{}

func (s *Server) ComputeDouble(ctx context.Context, request *impl.ComputeDoubleRequest, version apiversion.Version) (*impl.ComputeDoubleResponse, error) {
	in := request.Input64
	out := 2 * in

	response := &impl.ComputeDoubleResponse{}

	if sign(in) != sign(out) {
		// overflow
		response.Overflow = true
	} else {
		response.Response = out
	}

	return response, nil
}

func sign(x int64) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
```
All the boilerplate code to:
 * add a named pipe to the server for each version of the API group, listening for each version's requests, and replying with each version's responses
 * convert versioned requests to internal representations
 * convert internal responses back to versioned responses
 * create clients to talk to the API group and its versions
is generated automatically using [gengo](https://github.com/kubernetes/gengo).

The only caveat is that when conversions cannot be made trivially (e.g. when fields from internal and versioned `struct`s have different types), API devs need to define conversion functions. They can do that by creating an (otherwise optional) `pkg/server/<api_group_name>/impl/<version>/conversion.go` file, containing functions of the form `func convert_pb_<Type>_To_impl_<Type>(in *pb.<Type>, out *impl.<Type>) error` or `func convert_impl_<Type>_To_pb_<Type>(in *impl.<Type>, out *pb.<Type>) error`; for example, in our `dummy` example above, we need to define a conversion function to account for the different fields in requests and responses from `v1alpha1` to `v1`; so `pkg/server/dummy/impl/v1alpha1/conversion.go` could look like:
```go
func convert_pb_ComputeDoubleRequest_To_impl_ComputeDoubleRequest(in *pb.ComputeDoubleRequest, out *impl.ComputeDoubleRequest) error {
	out.Input64 = int64(in.Input32)
	return nil
}

func convert_impl_ComputeDoubleResponse_To_pb_ComputeDoubleResponse(in *impl.ComputeDoubleResponse, out *pb.ComputeDoubleResponse) error {
	i := in.Response
	if i > math.MaxInt32 || i < math.MinInt32 {
		return fmt.Errorf("int32 overflow for %d", i)
	}
	out.Response32 = int32(i)
	return nil
}
```

## How to change the API

Existing API versions are immutable.

### How to add a new API group

Simply create a new `client/api/<api_group_name>/<version>/api.proto` file, defining your new service; then generate the Go protobuf code, and run the CSI-proxy generator to generate all the boilerplate code.

FIXME: add more details on which commands to run, and which files to edit when done generating.

### How to add a new version to an existing API group

Any changes to the API of an existing API group requires creating a new API version.

Steps to add a new API version:
1. define it its own new `api.proto` file
2. generate the Go protobuf code
3. update the API group's internal representations (in its `types.go` file) to be able to represent all of the group's version (the new and the old ones)
4. add any needed conversion functions for all existing versions of this API group to account for the changes made at the previous step
5. re-generate all of the Go boilerplate code
6. now you can change the API group's server to add your new feature!

### How to deprecate, and eventually remove...

From the CSI [proxy KEP](https://github.com/kubernetes/enhancements/blob/master/keps/sig-windows/20190714-windows-csi-support.md#csi-proxy-grpc-api-graduation-and-deprecation-policy):

> In accordance with standard Kubernetes conventions, the above API will be introduced as v1alpha1 and graduate to v1beta1 and v1 as the feature graduates. Beyond a vN release in the future, new RPCs and enhancements to parameters will be introduced through vN+1alpha1 and graduate to vN+1beta1 and vN+1 stable versions as the new APIs mature.
>
> Members of CSIProxyService API may be deprecated and then removed from csi-proxy.exe in a manner similar to Kubernetes deprecation [policy](https://kubernetes.io/docs/reference/using-api/deprecation-policy/) although maintainers will make an effort to ensure such deprecation is as rare as possible. After their announced deprecation, a member of CSIProxyService API must be supported:
>
> 1. 12 months or 3 releases (whichever is longer) if the API member is part of a Stable/vN version.
> 2. 9 months or 3 releases (whichever is longer) if the API member is part of a Beta/vNbeta1 version.
> 3. 0 releases if the API member is part of an Alpha/vNalpha1 version.

With that in mind, each subsection below details how to deprecate, then remove:

#### A field in an API object

Mark it as deprecated in the `proto` file, e.g.:
```proto
message ComputeDoubleRequest{
    int32 input32 = 1 [deprecated=true];
}
```
then regenerate the protobuf code.

For removal, simply remove it in the first API version that doesn't support that field any more.

#### An API procedure

Similarly, mark the procedure as deprecated in the protobuf definition, e.g.:
```proto
service Dummy {
    // ComputeDouble computes the double of the input. Real smart stuff!
    rpc ComputeDouble(ComputeDoubleRequest) returns (ComputeDoubleResponse) {
        option deprecated = true;
    }
}
```
then regenerate the protobuf code, and remove it from the first API version that doesn't support that procedure any more.

#### An API version

Again, mark the version as deprecated in its protobuf definition, e.g.:
```proto
// Deprecated: Do not use.
// v1alpha1 is no longer maintained, and will be removed soon.
package v1alpha1;

service Dummy {
    ...
}
```
then regenerate the protobuf code, and run `csi-proxy-api-gen`: it will also mark the version's server (``) and client (``) packages as deprecated.

For removal, remove the whole `client/api/<api_group_name>/<version>` directory, and run `csi-proxy-api-gen`, it will remove all references to the removed version.

#### An API group

Deprecate and remove all its versions as explained in the previous version; then remove the entire `client/api/<api_group_name>` directory, and run `csi-proxy-api-gen`, it will remove all references to the removed API group.

## Detailed breakdown of generated files

This section details how `csi-proxy-api-gen` works, and what files it generates; `csi-proxy-api-gen` is built on top of [gengo](https://github.com/kubernetes/gengo), and re-uses part of [k8s' code-generator](https://github.com/kubernetes/code-generator), notably to generate conversion functions.

First, it looks for all API group definitions, which are either subdirectories of `client/api/`, or any go package whose `doc.go` file contains a `// +csi-proxy-api-gen` comment.

Then for each API group it finds:
1. it iterates through each version subpackage, and in each looks for the `<ApiGroupName>Server` interface, and compiles the list of callbacks that the group's `Server` needs to implement as well as the list of top-level `struct`s (`*Request`s and `*Response`s)
2. it looks for an existing `pkg/server/<api_group_name>/impl/types.go` file:
    * if it exists, it checks that it contains all the expected top-level `struct`s from the previous step
    * if it doesn't exist, _and_ the API group only defines one version, it auto-generates one that simply copies the protobuf `struct`s (from the previous step) - this is meant to make it easy to bootstrap a new API group
3. it generates the `pkg/server/<api_group_name>/impl/types_generated.go` file, using the list of callbacks from the first step above
4. if `pkg/server/<api_group_name>/server.go` doesn't exist, it generates a skeleton for it - this, too, is meant to make it easy to bootstrap new API groups
5. then for each version of the API:
    1. it looks for an existing `pkg/server/<api_group_name>/impl/<version>/conversion.go`, generates an empty one if it doesn't exist; then looks for existing conversion functions
    2. it generates missing conversion functions to `pkg/server/<api_group_name>/impl/<version>/conversion_generated.go`
    3. it generates `pkg/server/<api_group_name>/impl/<version>/server_generated.go`
6. it generates `pkg/server/<api_group_name>/impl/api_group_generated.go` to list all the versioned servers it's just created
7. and finally, it generates `client/groups/<api_group_name>/<version>/client_generated.go` for each version

When `csi-proxy-api-gen` has successfully run to completion, [our example API group's go package from earlier](#serverPkgTree) will look something like:
```
pkg/server/<api_group_name>
├── api_group_generated.go
├── impl
│   ├── types.go
│   ├── types_generated.go
│   ├── v1
│   │   ├── conversion.go
│   │   ├── conversion_generated.go
│   │   └── server_generated.go
│   └── v1alpha1
│       ├── conversion.go
│       ├── conversion_generated.go
│       └── server_generated.go
└── server.go
```
