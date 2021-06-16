package volume

import (
	"context"
	"fmt"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/os/volume"
	internal "github.com/kubernetes-csi/csi-proxy/pkg/server/volume/impl"
)

type fakeVolumeAPI struct {
	diskVolMap map[uint32][]string
}

var _ volume.API = &fakeVolumeAPI{}

func (volumeAPI *fakeVolumeAPI) Fill(diskToVolMapIn map[uint32][]string) {
	for d, v := range diskToVolMapIn {
		volumeAPI.diskVolMap[d] = v
	}
}

func (volumeAPI *fakeVolumeAPI) ListVolumesOnDisk(diskNumber uint32, partitionNumber uint32) (volumeIDs []string, err error) {
	v := volumeAPI.diskVolMap[diskNumber]
	if v == nil {
		return nil, fmt.Errorf("returning error for %d list", diskNumber)
	}
	return v, nil
}

func (volumeAPI *fakeVolumeAPI) MountVolume(volumeID, path string) error {
	return nil
}

func (volumeAPI *fakeVolumeAPI) UnmountVolume(volumeID, path string) error {
	return nil
}

func (volumeAPI *fakeVolumeAPI) IsVolumeFormatted(volumeID string) (bool, error) {
	return true, nil
}

func (volumeAPI *fakeVolumeAPI) FormatVolume(volumeID string) error {
	return nil
}

func (volumeAPI *fakeVolumeAPI) ResizeVolume(volumeID string, size int64) error {
	return nil
}

func (volumeAPI *fakeVolumeAPI) GetDiskNumberFromVolumeID(volumeID string) (uint32, error) {
	return 0, nil
}

func (volumeAPI *fakeVolumeAPI) GetVolumeIDFromTargetPath(mount string) (string, error) {
	return "id", nil
}

func (volumeAPI *fakeVolumeAPI) GetVolumeStats(volumeID string) (int64, int64, error) {
	return -1, -1, nil
}

func (volumeAPI *fakeVolumeAPI) WriteVolumeCache(volumeID string) error {
	return nil
}

func TestListVolumesOnDisk(t *testing.T) {
	v1, err := apiversion.NewVersion("v1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}

	testCases := []struct {
		name              string
		inputDiskNumber   uint32
		expectedVolumeIds []string
		isErrorExpected   bool
		expectedError     error
	}{
		{
			name:              "return two volumeIDs",
			inputDiskNumber:   1,
			expectedVolumeIds: []string{"volumeID1", "volumeID2"},
			isErrorExpected:   false,
			expectedError:     nil,
		},
		{
			name:              "return one volumeIDs",
			inputDiskNumber:   2,
			expectedVolumeIds: []string{"volumeID3"},
			isErrorExpected:   false,
			expectedError:     nil,
		},
		{
			name:              "return error",
			inputDiskNumber:   3,
			expectedVolumeIds: nil,
			isErrorExpected:   true,
			expectedError:     fmt.Errorf("returning error for 3 list"),
		},
	}

	diskToVolMap := map[uint32][]string{
		1: {"volumeID1", "volumeID2"},
		2: {"volumeID3"},
	}
	volAPI := &fakeVolumeAPI{
		diskVolMap: make(map[uint32][]string),
	}
	volAPI.Fill(diskToVolMap)

	volumeSrv, err := NewServer(volAPI)
	if err != nil {
		t.Fatalf("Volume server could not be initialized: %v", err)
	}

	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		listInput := &internal.ListVolumesOnDiskRequest{
			DiskNumber: tc.inputDiskNumber,
		}
		volumeListResponse, err := volumeSrv.ListVolumesOnDisk(context.TODO(), listInput, v1)
		if tc.isErrorExpected {
			if tc.expectedError.Error() != err.Error() {
				t.Fatalf("Expected error: %v. Got error: %v", tc.expectedError, err)
			}
		} else {
			if err != nil {
				t.Fatalf("Error %v not expected", err)
			}

			expectedVolumeIDMap := make(map[string]int)
			for _, j := range tc.expectedVolumeIds {
				expectedVolumeIDMap[j] = 0
			}
			for _, i := range volumeListResponse.VolumeIds {
				if _, found := expectedVolumeIDMap[i]; found == true {
					expectedVolumeIDMap[i]++
				} else {
					t.Fatalf("Found unexpected volume: %s", i)
				}
			}
			for k, v := range expectedVolumeIDMap {
				if v != 1 {
					t.Fatalf("Volume: %s count: %d", k, v)
				}
			}
		}
	}
}
