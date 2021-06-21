# Changelog since v0.2.2 (Beta.2)

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- `--kubelet-pod-path` and `--kubelet-csi-plugins-path` parameters for csi-proxy.exe is replaced with `--kubelet-path`. No changes are necessary if default values for the parameters were being used to initialize csi-proxy.exe
- All CSI Plugins (and other clients of csi-proxy) using v1alpha and v1beta named pipes and APIs should continue to function as is.

## Changes by Kind

### Feature

- Combine plugin and pod paths into one KubeletPath ([#150](https://github.com/kubernetes-csi/csi-proxy/pull/150), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Make GPT the new partition style for volumes ([#128](https://github.com/kubernetes-csi/csi-proxy/pull/128), [@manueltellez](https://github.com/manueltellez))
- Add NewClientWithPipePath per API that can communicate with a custom named pipe ([#124](https://github.com/kubernetes-csi/csi-proxy/pull/124), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Return disk serial number in ListDiskIDs ([#116](https://github.com/kubernetes-csi/csi-proxy/pull/116), [@wongma7](https://github.com/wongma7))
- Update klog to version 2 ([#114](https://github.com/kubernetes-csi/csi-proxy/pull/114), [@jingxu97](https://github.com/jingxu97))
- Support reporting Disk IDs with CodeSet StorageIDCodeSetBinary ([#110](https://github.com/kubernetes-csi/csi-proxy/pull/110), [@gab-satchi](https://github.com/gab-satchi))

### API
- v1beta2 SMB API group with error field removed from all API response msgs ([#146](https://github.com/kubernetes-csi/csi-proxy/pull/146), [@mauriciopoppe](https://github.com/mauriciopoppe))
- v1beta2 FileSystem API group with renamed APIs LinkPath -> CreateSymlink and IsMountPoint -> IsSymlink and error field removed from all API response msgs ([#143](https://github.com/kubernetes-csi/csi-proxy/pull/143), [@mauriciopoppe](https://github.com/mauriciopoppe))
- v1beta3 Disk API group with renamed APIs DiskStats -> GetDiskStats, SetAttachState -> SetDiskState and SetAttachState -> GetDiskState ([#140](https://github.com/kubernetes-csi/csi-proxy/pull/140), [@mauriciopoppe](https://github.com/mauriciopoppe))
- v1beta3 Volume API group with renamed APIs DismountVolume -> UnmountVolume, VolumeStats -> GetVolumeStats, GetVolumeDiskNumber -> GetDiskNumberFromVolumeID and GetVolumeIDFromTargetPath -> GetVolumeIDFromMount ([#138](https://github.com/kubernetes-csi/csi-proxy/pull/138), [@mauriciopoppe](https://github.com/mauriciopoppe))
- v1alpha2 iSCSI API group with new API SetMutualChapSecret ([#102](https://github.com/kubernetes-csi/csi-proxy/pull/102), [@jmpfar](https://github.com/jmpfar))

### Bug or Regression
- Close Open disk file handles ([#147](https://github.com/kubernetes-csi/csi-proxy/pull/147), [@wongma7](https://github.com/wongma7))
- Don't count reserved partitions (gpt) when checking if any exist ([#145](https://github.com/kubernetes-csi/csi-proxy/pull/145), [@wongma7](https://github.com/wongma7))
- Normalize windows paths in SMB apis to support linux style paths in source ([#128](https://github.com/kubernetes-csi/csi-proxy/pull/128), [@manueltellez](https://github.com/manueltellez))
- Fix IsVolumeFormatted issue on checking volume type "Unknown" ([#127](https://github.com/kubernetes-csi/csi-proxy/pull/127), [@jingxu97](https://github.com/jingxu97))
- Remove specific check for NTFS in IsVolumeFormatted ([#123](https://github.com/kubernetes-csi/csi-proxy/pull/123), [@jingxu97](https://github.com/jingxu97))
- Fix smb mount PermissionDenied issue on Windows ([#117](https://github.com/kubernetes-csi/csi-proxy/pull/117), [@andyzhangx](https://github.com/andyzhangx))

### Other (Cleanup or Flake)
- Rename /internal to pkg and nested /internal/server/<group>/internal to impl ([#152](https://github.com/kubernetes-csi/csi-proxy/pull/152), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Remove redundant alias import and replace few function names in ISCSI server ([#137](https://github.com/kubernetes-csi/csi-proxy/pull/137), [@humblec](https://github.com/humblec))
- Replace and correct imports in integration tests ([#136](https://github.com/kubernetes-csi/csi-proxy/pull/136), [@humblec](https://github.com/humblec))
- Add asciiflow diagram to csi-proxy ([#106](https://github.com/kubernetes-csi/csi-proxy/pull/106), [@jayunit100](https://github.com/jayunit100))
