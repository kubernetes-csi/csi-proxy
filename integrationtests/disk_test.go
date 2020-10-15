package integrationtests

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta2"
	v1beta2client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1beta2"
	"github.com/stretchr/testify/require"
)

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
func TestDiskAPIGroupV1Beta1(t *testing.T) {
	t.Run("ListDiskIDs", func(t *testing.T) {
		client, err := v1beta2client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		diskNumber := 0
		id := "page83"
		listRequest := &v1beta2.ListDiskIDsRequest{}
		diskIDsResponse, err := client.ListDiskIDs(context.TODO(), listRequest)
		require.Nil(t, err)

		cmd := "hostname"
		hostname, err := runPowershellCmd(cmd)
		if err != nil {
			t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, hostname)
		}

		hostname = strings.TrimSpace(hostname)

		diskIDsMap := diskIDsResponse.GetDiskIDs()
		if diskIDsMap != nil {
			diskIDs, found := diskIDsMap[strconv.Itoa(diskNumber)]
			if !found {
				t.Errorf("Cannot find Disk %d", diskNumber)
			}

			idValue, found := diskIDs.Identifiers[id]
			if !found {
				t.Errorf("Cannot find %s ID of Disk %d", id, diskNumber)
			}

			// In GCE, the page83 ID of Disk 0 contains the hostname
			if !strings.Contains(idValue, hostname) {
				t.Errorf("%s ID of Disk %d is incorrect. Expected to contain: %s. Received: %s", id, diskNumber, hostname, idValue)
			}
		}
	})
}
