# Changelog since v0.1.0 (Alpha)

## Urgent Upgrade Notes 

### (No, really, you MUST read this before you upgrade)

- None. All clients using v1alpha named pipes should continue to function as is.

## Changes by Kind

### Feature

- ListDiskIDs API in v1beta1 Disk API group ([#61](https://github.com/kubernetes-csi/csi-proxy/pull/61), [@ksubrmnn](https://github.com/ksubrmnn))
- DiskStats API in v1beta1 Disk API group ([#65](https://github.com/kubernetes-csi/csi-proxy/pull/65), [@manueltellez](https://github.com/manueltellez))
- VolumeStats, GetVolumeDiskNumber and GetVolumeIDFromMount API in v1beta1 Volume API group ([#65](https://github.com/kubernetes-csi/csi-proxy/pull/65), [@manueltellez](https://github.com/manueltellez))
- Run csi-proxy as a Windows service ([#62](https://github.com/kubernetes-csi/csi-proxy/pull/62), [@ddebroy](https://github.com/ddebroy))
- Retrieve SCSI Page 83 IDs directly using IOCTL_STORAGE_QUERY_PROPERTY ([#42](https://github.com/kubernetes-csi/csi-proxy/pull/52), [@ksubrmnn](https://github.com/ksubrmnn))

### Bug or Regression
- Test SMB path validity during remounts ([#66](https://github.com/kubernetes-csi/csi-proxy/pull/66), [@andyzhangx](https://github.com/andyzhangx))

### Other (Cleanup or Flake)
- Build and publish csi-proxy binary to gcs bucket k8s-artifacts-csi ([#53](https://github.com/kubernetes-csi/csi-proxy/pull/53), [@jingxu97](https://github.com/jingxu97))
