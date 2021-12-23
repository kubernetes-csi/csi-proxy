// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/kubernetes-csi/csi-proxy/client"
	v1 "github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1"
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"google.golang.org/grpc"
)

// GroupName is the group name of this API.
const GroupName = "filesystem"

// Version is the api version.
var Version = apiversion.NewVersionOrPanic("v1")

type Client struct {
	client     v1.FilesystemClient
	connection *grpc.ClientConn
}

// NewClient returns a client to make calls to the filesystem API group version v1.
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

	client := v1.NewFilesystemClient(connection)
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
var _ v1.FilesystemClient = &Client{}

func (w *Client) CreateSymlink(context context.Context, request *v1.CreateSymlinkRequest, opts ...grpc.CallOption) (*v1.CreateSymlinkResponse, error) {
	return w.client.CreateSymlink(context, request, opts...)
}

func (w *Client) IsSymlink(context context.Context, request *v1.IsSymlinkRequest, opts ...grpc.CallOption) (*v1.IsSymlinkResponse, error) {
	return w.client.IsSymlink(context, request, opts...)
}

func (w *Client) Lsdir(context context.Context, request *v1.LsdirRequest, opts ...grpc.CallOption) (*v1.LsdirResponse, error) {
	return w.client.Lsdir(context, request, opts...)
}

func (w *Client) Mkdir(context context.Context, request *v1.MkdirRequest, opts ...grpc.CallOption) (*v1.MkdirResponse, error) {
	return w.client.Mkdir(context, request, opts...)
}

func (w *Client) PathExists(context context.Context, request *v1.PathExistsRequest, opts ...grpc.CallOption) (*v1.PathExistsResponse, error) {
	return w.client.PathExists(context, request, opts...)
}

func (w *Client) Rmdir(context context.Context, request *v1.RmdirRequest, opts ...grpc.CallOption) (*v1.RmdirResponse, error) {
	return w.client.Rmdir(context, request, opts...)
}
