# Changelog since v0.2.1 (Beta.1)

## Urgent Upgrade Notes 

### (No, really, you MUST read this before you upgrade)

- None. All clients using v1alpha and v1beta named pipes should continue to function as is.

## Changes by Kind

### Feature

- v1beta2 Volume API group with new API WriteVolumeCache ([#86](https://github.com/kubernetes-csi/csi-proxy/pull/86), [@jingxu97](https://github.com/jingxu97))
- v1alpha1 System API group with new API GetBIOSSerialNumber ([#81](https://github.com/kubernetes-csi/csi-proxy/pull/81), [@ksubrmnn](https://github.com/ksubrmnn))
- v1beta2 Disk API group with new APIs SetAttachState and GetAttachState ([#95](https://github.com/kubernetes-csi/csi-proxy/pull/95), [@jmpfar](https://github.com/jmpfar))
- v1alpha1 System API group with new APIs StartService StopService and GetService ([#100](https://github.com/kubernetes-csi/csi-proxy/pull/100), [@jmpfar](https://github.com/jmpfar))
- v1alpha1 iSCSI API group with new APIs AddTargetPortal DiscoverTargetPortal RemoveTargetPortal ListTargetPortals ConnectTarget DisconnectTarget and GetTargetDisks ([#99](https://github.com/kubernetes-csi/csi-proxy/pull/99), [@jmpfar](https://github.com/jmpfar))

### Bug or Regression
- Simplify lookup of symlink targets ([#96](https://github.com/kubernetes-csi/csi-proxy/pull/96), [@jingxu97](https://github.com/jingxu97))
- Add trailing backslash to remote path when linking to SMB shares ([#98](https://github.com/kubernetes-csi/csi-proxy/pull/98), [@marosset](https://github.com/marosset))
- No-op when resize is called and size of partition is >= requested size ([#91](https://github.com/kubernetes-csi/csi-proxy/pull/91), [@manueltellez ](https://github.com/manueltellez))

### Other (Cleanup or Flake)
- Add vendor directory ([#82](https://github.com/kubernetes-csi/csi-proxy/pull/82), [@jingxu97](https://github.com/jingxu97))
- Updates to README covering installation of csi-proxy as a service ([#86](https://github.com/kubernetes-csi/csi-proxy/pull/86), [@jingxu97](https://github.com/jingxu97))
- Bump go version to 1.13 in go.mod ([#89](https://github.com/kubernetes-csi/csi-proxy/pull/89), [@mayankshah1607](https://github.com/mayankshah1607))
- Enable integration tests on Windows through Github Actions ([#90](https://github.com/kubernetes-csi/csi-proxy/pull/90), [@mayankshah1607](https://github.com/mayankshah1607))
