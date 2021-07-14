# Changelog since v1.0.0-rc1 (v1 RC1)

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- None

## Changes by Kind

### Feature

- None

### API
- No changes since v1 RC1
- Latest versions supported:
  - [Disk v1](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/disk/v1/api.proto)
  - [Volume v1](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/volume/v1/api.proto)
  - [FileSystem v1](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/filesystem/v1/api.proto)
  - [SMB v1](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/smb/v1/api.proto)
  - [iSCSI v1Alpha2](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/iscsi/v1alpha2/api.proto)
  - [system v1Alpha1](https://github.com/kubernetes-csi/csi-proxy/blob/master/client/api/system/v1alpha1/api.proto)

### Bug or Regression
- None

### Other (Cleanup or Flake)
- Add logic to upload binary to gcs with version information ([#148](https://github.com/kubernetes-csi/csi-proxy/pull/148), [@jingxu97](https://github.com/jingxu97))
- Update node-driver-registrar image tag (for Windows) in README.md ([#162](https://github.com/kubernetes-csi/csi-proxy/pull/162), [@mauriciopoppe](https://github.com/mauriciopoppe))
