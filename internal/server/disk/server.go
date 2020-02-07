package disk

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

type Server struct {
	hostAPI API
}

type API interface {
	IsDiskInitialized(diskID string) (bool, error)
	InitializeDisk(diskID string) error
	PartitionsExist(diskID string) (bool, error)
	CreatePartition(diskID string) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) ListDiskLocations(context context.Context, request *internal.ListDiskLocationsRequest, version apiversion.Version) (*internal.ListDiskLocationsResponse, error) {
	// TODO: auto-generated stub
	return nil, nil
}

func (s *Server) PartitionDisk(context context.Context, request *internal.PartitionDiskRequest, version apiversion.Version) (*internal.PartitionDiskResponse, error) {
	response := &internal.PartitionDiskResponse{}
	diskID := request.DiskId

	initialized, err := s.hostAPI.IsDiskInitialized(diskID)
	if err != nil {
		return response, err
	}
	if !initialized {
		err = s.hostAPI.InitializeDisk(diskID)
		if err != nil {
			return response, err
		}
	}

	paritioned, err := s.hostAPI.PartitionsExist(diskID)
	if err != nil {
		return response, err
	}
	if !paritioned {
		err = s.hostAPI.CreatePartition(diskID)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}
