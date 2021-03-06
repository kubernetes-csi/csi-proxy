// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/kubernetes-csi/csi-proxy/client"
	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"google.golang.org/grpc"
)

// GroupName is the group name of this API.
const GroupName = "disk"

// Version is the api version.
var Version = apiversion.NewVersionOrPanic("v1alpha1")

type Client struct {
	client     v1alpha1.DiskClient
	connection *grpc.ClientConn
}

// NewClient returns a client to make calls to the disk API group version v1alpha1.
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

	client := v1alpha1.NewDiskClient(connection)
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
var _ v1alpha1.DiskClient = &Client{}

func (w *Client) GetDiskNumberByName(context context.Context, request *v1alpha1.GetDiskNumberByNameRequest, opts ...grpc.CallOption) (*v1alpha1.GetDiskNumberByNameResponse, error) {
	return w.client.GetDiskNumberByName(context, request, opts...)
}

func (w *Client) ListDiskLocations(context context.Context, request *v1alpha1.ListDiskLocationsRequest, opts ...grpc.CallOption) (*v1alpha1.ListDiskLocationsResponse, error) {
	return w.client.ListDiskLocations(context, request, opts...)
}

func (w *Client) PartitionDisk(context context.Context, request *v1alpha1.PartitionDiskRequest, opts ...grpc.CallOption) (*v1alpha1.PartitionDiskResponse, error) {
	return w.client.PartitionDisk(context, request, opts...)
}

func (w *Client) Rescan(context context.Context, request *v1alpha1.RescanRequest, opts ...grpc.CallOption) (*v1alpha1.RescanResponse, error) {
	return w.client.Rescan(context, request, opts...)
}
