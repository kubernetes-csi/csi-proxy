# Usage

## Overview

CSI Proxy organizes the functionalities it provides into different API groups, which are:

- [Disk](/pkg/disk/)
- [FileSystem](/pkg/filesystem/)
- [iSCSI](/pkg/iscsi/) (experimental)
- [SMB](/pkg/smb/)
- [System](/pkg/system/) (experimental)
- [Volume](/pkg/volume/)

Each API group defines an interface specifying its API and provides a struct implementing the interface.
The API-level interface takes in a `HostAPI` that handles the implementation details and is only responsible for input/output checking.
The user-facing interface is exposed in `/pkg/<API group>/`, whereas the `HostAPI` for each API group is exposed via `/pkg/<API group>/hostapi/`.
Both these paths expose a method `New` that returns an instance of the API interface and `HostAPI`, respectively.
The relevant request/response types for the APIs are also provided under `/pkg/<API Group>`.

## Usage

### Go Code

To use any API group, the driver needs to import the API group and its `HostAPI`.
Below is an example for using the `FileSystem` API group.

```go
import fsapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem/hostapi"
import fs "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem"

func NewCSIProxyMounterV1() (*CSIProxyMounterV1, error) {
    fsClient, err := fs.New(fsapi.New())
    if err != nil {
        return nil, err
    }
    return &CSIProxyMounterV1{
        FsClient:     fsClient,
    }, nil
}

func (mounter *CSIProxyMounterV1) PathExists(path string) (bool, error) {
    isExistsResponse, err := mounter.FsClient.PathExists(context.Background(),
        // the request type is exposed in pkg/filesystem
        &fs.PathExistsRequest{
            Path: mount.NormalizeWindowsPath(path),
        })
    if err != nil {
        return false, err
    }
    return isExistsResponse.Exists, err
}
```

### Deployment

CSI driver containers need to run as HostProcess pods for the CSI Proxy commands to complete privileged operations.
This can be done by updating the driver pod spec, typically embedded in a Daemonset.
The Kubernetes Windows HostProcess Pod [docs](https://kubernetes.io/docs/tasks/configure-pod-container/create-hostprocess-pod/) goes into more detail about each field.

```yaml
spec:
  securityContext:
    hostNetwork: true
    windowsOptions:
      hostProcess: true
      runAsUserName: "NT AUTHORITY\\SYSTEM"
```

Using HostProcess pods have a few important consequences on the deployed containers, as noted in the HostProcess docs.
- HostProcess containers have no file system or resource isolation, so they have a complete view of the host machine’s file system.
This means that paths passed in **paths passed as command line arguments must be absolute paths with respect to the host**.
On the other hand, depending on the containerd and the Windows OS version, Kubernetes volume mounts are either mounted relative to a subdirectory of the host machine specified by an environment variable `%CONTAINER_SANDBOX_MOUNT_POINT%` or **mounted relative to the host process root**. See [HostProcess Caveats](#hostprocess-caveats).
- HostProcess pods can only contain HostProcess containers.
Often, the CSI node registrar is deployed in the same pod as the driver, so file paths for both containers need to be updated.
- Named pipes and Unix domain sockets are not supported and should be accessed via their absolute path with respect to the host.

Instead of mounting host process paths such as `c:\var\lib\kubelet` to each container’s own filesystem, the container can directly access these paths.

If drivers are migrating from CSI Proxy v1, note that since the current version (v2) no longer has a separate binary running on the host, drivers no longer need to mount named pipes for each API group/version used. Related volumes and volume mounts can be safely deleted.

Here is an example driver's deployment configuration.

```yaml
kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: csi-storage-node-win
spec:
  selector:
    matchLabels:
      app: csi-driver-win
  template:
    metadata:
      labels:
        app: csi-driver-win
    spec:
      serviceAccountName: csi-node-sa
      tolerations:
      - key: "node.kubernetes.io/os"
        operator: "Equal"
        value: "win1809"
        effect: "NoSchedule"
      nodeSelector:
        kubernetes.io/os: windows
      containers:
        - name: csi-driver-registrar
          image: k8s.gcr.io/sig-storage/csi-node-driver-registrar:v2.1.0
          args:
            - --v=5
            - --csi-address=unix://c:\var\lib\kubelet\plugins\csi.org.io\csi.sock
            - --kubelet-registration-path=/var/lib/kubelet/plugins/csi.org.io/csi.sock
            - --plugin-registration-path=/var/lib/kubelet/plugins_registry
          env:
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
        - name: csi-driver
          # placeholder, use your CSI driver
          image: org/csi-driver:win-v1
          args:
            - --v=5
            - --endpoint=unix:/c:\var\lib\kubelet\plugins\csi.org.io\csi.sock
      volumes:
        - name: plugin-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi.org.io
            type: DirectoryOrCreate
```

### HostProcess Caveats

HostProcess containers require containerd v1.6 to work, but new Windows OS APIs available in Windows Server 2019 starting from July 2022 allow for a cleaner implementation of HostProcess containers using *bind mount* volume mount behaviors. The new behavior requires containerd v1.7, which is not yet released at the time of writing. Running containerd v1.7 with an older version of Windows not supporting the new APIs would cause HostProcess containers to fail.

Practically, the difference is that HostProcess containers running on nodes with containerd v1.6 see the whole host file system, whereas HostProcess containers running on nodes with containerd v1.7 are presented with a merged view of the host OS’s file system and container-local volumes. Containerd v1.6 mounts the container files at `c:\C\{container-id}`, whereas containerd v1.7 mounts the container files at `c:\hpc`, where changes inside that path are only visible to the container. In both versions, the mount path is exposed as an environment variable `%CONTAINER_SANDBOX_MOUNT_POINT%`, which is also set as the default working directory of the containers.

Another key difference is that volume mounts in containerd v1.6 are relative to the container mount path (i.e., `%CONTAINER_SANDBOX_MOUNT_POINT%`), whereas volume mounts in containerd v1.7 are relative to the root of the host file system. Drivers that rely on mounting volumes to containers are likely going to be broken by containerd version changes. Therefore, it’s recommended to migrate drivers deployment specs to not use any volume mounts and instead rely on absolute host file system paths. This should ensure that the driver will work with both containerd v1.6 and v1.7.
