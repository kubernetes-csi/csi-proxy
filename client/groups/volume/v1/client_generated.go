// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/kubernetes-csi/csi-proxy/client"
	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1"
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"google.golang.org/grpc"
)

// GroupName is the group name of this API.
const GroupName = "volume"

// Version is the api version.
var Version = apiversion.NewVersionOrPanic("v1")

type Client struct {
	client     v1.VolumeClient
	connection *grpc.ClientConn
}

// NewClient returns a client to make calls to the volume API group version v1.
// It's the caller's responsibility to Close the client when done.
func NewClient() (*Client, error) {
	pipePath := client.PipePath(GroupName, Version)
	return NewClientWithPipePath(pipePath)
}

// NewClientWithPipePath returns a client to make calls to the named pipe located at "pipePath".
// It's the caller's responsibility to Close the client when done.
func NewClientWithPipePath(pipePath string) (*Client, error) {

	// verify that the pipe exists
	_, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		return nil, err
	}

	connection, err := grpc.Dial(pipePath,
		grpc.WithContextDialer(func(context context.Context, s string) (net.Conn, error) {
			return winio.DialPipeContext(context, s)
		}),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := v1.NewVolumeClient(connection)
	return &Client{
		client:     client,
		connection: connection,
	}, nil
}

// Close closes the client. It must be called before the client gets GC-ed.
func (w *Client) Close() error {
	return w.connection.Close()
}

// ensures we implement all the required methods
var _ v1.VolumeClient = &Client{}

func (w *Client) FormatVolume(context context.Context, request *v1.FormatVolumeRequest, opts ...grpc.CallOption) (*v1.FormatVolumeResponse, error) {
	return w.client.FormatVolume(context, request, opts...)
}

func (w *Client) GetClosestVolumeIDFromTargetPath(context context.Context, request *v1.GetClosestVolumeIDFromTargetPathRequest, opts ...grpc.CallOption) (*v1.GetClosestVolumeIDFromTargetPathResponse, error) {
	return w.client.GetClosestVolumeIDFromTargetPath(context, request, opts...)
}

func (w *Client) GetDiskNumberFromVolumeID(context context.Context, request *v1.GetDiskNumberFromVolumeIDRequest, opts ...grpc.CallOption) (*v1.GetDiskNumberFromVolumeIDResponse, error) {
	return w.client.GetDiskNumberFromVolumeID(context, request, opts...)
}

func (w *Client) GetVolumeIDFromTargetPath(context context.Context, request *v1.GetVolumeIDFromTargetPathRequest, opts ...grpc.CallOption) (*v1.GetVolumeIDFromTargetPathResponse, error) {
	return w.client.GetVolumeIDFromTargetPath(context, request, opts...)
}

func (w *Client) GetVolumeStats(context context.Context, request *v1.GetVolumeStatsRequest, opts ...grpc.CallOption) (*v1.GetVolumeStatsResponse, error) {
	return w.client.GetVolumeStats(context, request, opts...)
}

func (w *Client) IsVolumeFormatted(context context.Context, request *v1.IsVolumeFormattedRequest, opts ...grpc.CallOption) (*v1.IsVolumeFormattedResponse, error) {
	return w.client.IsVolumeFormatted(context, request, opts...)
}

func (w *Client) ListVolumesOnDisk(context context.Context, request *v1.ListVolumesOnDiskRequest, opts ...grpc.CallOption) (*v1.ListVolumesOnDiskResponse, error) {
	return w.client.ListVolumesOnDisk(context, request, opts...)
}

func (w *Client) MountVolume(context context.Context, request *v1.MountVolumeRequest, opts ...grpc.CallOption) (*v1.MountVolumeResponse, error) {
	return w.client.MountVolume(context, request, opts...)
}

func (w *Client) ResizeVolume(context context.Context, request *v1.ResizeVolumeRequest, opts ...grpc.CallOption) (*v1.ResizeVolumeResponse, error) {
	return w.client.ResizeVolume(context, request, opts...)
}

func (w *Client) UnmountVolume(context context.Context, request *v1.UnmountVolumeRequest, opts ...grpc.CallOption) (*v1.UnmountVolumeResponse, error) {
	return w.client.UnmountVolume(context, request, opts...)
}

func (w *Client) WriteVolumeCache(context context.Context, request *v1.WriteVolumeCacheRequest, opts ...grpc.CallOption) (*v1.WriteVolumeCacheResponse, error) {
	return w.client.WriteVolumeCache(context, request, opts...)
}
