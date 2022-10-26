# CSI Proxy

CSI Proxy is a Go library providing convenience methods to execute privileged storage operations in Windows, such as formatting and mounting volumes.
A container, such as CSI node plugins, can import the CSI Proxy library to get a Go interface for storage-related Windows system calls.
Since the commands executed are privileged instructions, containers must run as [HostProcess containers](https://kubernetes.io/docs/tasks/configure-pod-container/create-hostprocess-pod/).

Closely related functionalities are bundled as API groups that target specific areas of storage. The available API groups are

- Disk
- Filesystem
- SMB
- Volume
- iSCSI (experimental)
- System (experimental)

## Compatibility

Recommended K8s Version: 1.23

## Usage

See [usage.md](/docs/API.md) for detailed usage instructions, as well as some notes on migrating from v1.

## Community, Discussion, Contribution, and Support

Check out [development.md](./docs/DEVELOPMENT.md) for instructions to set up a development environment to run CSI Proxy.

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/csi-windows)
- [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-storage)

## Supported CSI Drivers

- [SMB CSI Driver](https://github.com/kubernetes-csi/csi-driver-smb/tree/master/deploy/example/windows). To see specifically how this driver is invoked, you can look at https://github.com/kubernetes-csi/csi-driver-smb/blob/master/pkg/mounter/safe_mounter_windows.go.

- [Azure Disk CSI Driver](https://github.com/kubernetes-sigs/azuredisk-csi-driver/tree/master/deploy/example/windows)

- [Azure File CSI Driver](https://github.com/kubernetes-sigs/azurefile-csi-driver/tree/master/deploy/example/windows).  See https://github.com/kubernetes-sigs/azurefile-csi-driver/blob/master/pkg/mounter/safe_mounter_windows.go as an example of the invocation path

- [Google Compute Engine Persistent Disk CSI Driver](https://github.com/kubernetes-sigs/gcp-compute-persistent-disk-csi-driver)

### Code of Conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
