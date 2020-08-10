package volume

import (
	"context"
	"fmt"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/volume/internal"
)

type fakeVolumeAPI struct {
	diskVolMap map[string][]string
}

func (volumeAPI *fakeVolumeAPI) Fill(diskToVolMapIn map[string][]string) {
	for d, v := range diskToVolMapIn {
		volumeAPI.diskVolMap[d] = v
	}
}

func (volumeAPI *fakeVolumeAPI) ListVolumesOnDisk(diskID string) (volumeIDs []string, err error) {
	v := volumeAPI.diskVolMap[diskID]
	if v == nil {
		return nil, fmt.Errorf("returning error for %s list", diskID)
	}
	return v, nil
}

func (volumeAPI *fakeVolumeAPI) MountVolume(volumeID, path string) error {
	return nil
}

func (volumeAPI *fakeVolumeAPI) DismountVolume(volumeID, path string) error {
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

func (volumeAPI *fakeVolumeAPI) VolumeStats(volumeID string) (int64, int64, int64, error) {
	return -1, -1, -1, nil
}

func TestListVolumesOnDisk(t *testing.T) {
	v1alpha1, err := apiversion.NewVersion("v1alpha1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}

	testCases := []struct {
		name              string
		inputDiskID       string
		expectedVolumeIds []string
		isErrorExpected   bool
		expectedError     error
	}{
		{
			name:              "return two volumeIDs",
			inputDiskID:       "diskID1",
			expectedVolumeIds: []string{"volumeID1", "volumeID2"},
			isErrorExpected:   false,
			expectedError:     nil,
		},
		{
			name:              "return one volumeIDs",
			inputDiskID:       "diskID2",
			expectedVolumeIds: []string{"volumeID3"},
			isErrorExpected:   false,
			expectedError:     nil,
		},
		{
			name:              "return error",
			inputDiskID:       "diskID3",
			expectedVolumeIds: nil,
			isErrorExpected:   true,
			expectedError:     fmt.Errorf("returning error for diskID3 list"),
		},
	}

	diskToVolMap := map[string][]string{
		"diskID1": {"volumeID1", "volumeID2"},
		"diskID2": {"volumeID3"},
	}
	volAPI := &fakeVolumeAPI{
		diskVolMap: make(map[string][]string),
	}
	volAPI.Fill(diskToVolMap)

	volumeSrv, err := NewServer(volAPI)
	if err != nil {
		t.Fatalf("Volume server could not be initialized: %v", err)
	}

	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		listInput := &internal.ListVolumesOnDiskRequest{
			DiskId: tc.inputDiskID,
		}
		volumeListResponse, err := volumeSrv.ListVolumesOnDisk(context.TODO(), listInput, v1alpha1)
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
