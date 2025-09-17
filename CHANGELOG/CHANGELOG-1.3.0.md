# Changelog since v1.2.1

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- None

## Changes by Kind

### Feature

* Use WMI to implement Volume API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/360
* Set up a metrics server with data served on /metrics by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/371
* Use WMI to implement System API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/375
* feat: Use WMI to implement Disk API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/376
* feat: Use WMI to implement iSCSI API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/377
* feat: Use WMI to implement Smb API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/378
* feat: Use associator to find WMI instances for volume APIs by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/388
* feat: Use WMI to create SMB Global Mapping to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/387
* feat: Use WMI to implement Service related System APIs to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/389
* feat: Migrate PathValid API to use Win32 API to reduce PowerShell overhead by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/390

### API

- None

### Bug or Regression

* Revert "Add min go runtime to be 1.23 and add  godebug winsymlink=0" by @mauriciopoppe in https://github.com/kubernetes-csi/csi-proxy/pull/369
* fix: Ensure COM threading apartment for API calls by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/392
* fix: Ensure IsSymlink works on Windows mounted folder by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/393

### Other (Cleanup or Flake)

* cleanup: refine resize volume error logging on Windows node by @andyzhangx in https://github.com/kubernetes-csi/csi-proxy/pull/383
* Add doc for WMI by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/384
* chore: Bump release-tools by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/395
* cleanup: Cross-port review comment from library-development branch to process branch by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/401
* Add laozc to OWNERS by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/402

## Dependencies

### Added
_Nothing has changed._

### Changed

* golang.org/x/sys: v0.28.0 → v0.32.0
* cleanup: Bump microsoft/wmi to 0.34.0 containing WMI method call fix by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/394
* Bump actions/checkout from 4 to 5 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/403

### Removed
_Nothing has changed._
