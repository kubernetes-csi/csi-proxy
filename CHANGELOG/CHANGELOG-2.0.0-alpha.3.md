# Changelog since v2.0.0-alpha.2

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- The library now requires Go 1.26.

## Changes by Kind

### Feature

- None

### API

* feat: add support for managing disk read-only state by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/438

### Bug or Regression

* Remove dependency on microsoft/wmi (library version) by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/416

### Other (Cleanup or Flake)

* fix: Pin GitHub action to exact SHA by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/417
* [library-development] Replace deprecated github.com/pkg/errors by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/460
* [library-development] chore: upgrade to go1.25 by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/461
* Add laozc to OWNERS by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/471

## Dependencies

### Added
_Nothing has changed._

### Changed

* github.com/stretchr/testify: v1.7.0 → v1.12.1
* golang.org/x/sys: v0.32.0 → v0.48.0 (fixes CVE-2026-39824)
* k8s.io/klog/v2: v2.9.0 → v2.140.0
* Bump actions/checkout from 6.0.2 to 7.0.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/448
* Bump actions/setup-go from 6.4.0 to 7.0.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/449
* Bump k8s.io/klog/v2 from 2.9.0 to 2.140.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/450
* Bump github.com/stretchr/testify from 1.7.0 to 1.11.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/455
* Update release-tools and dependencies in library-development by @torredil in https://github.com/kubernetes-csi/csi-proxy/pull/479

### Removed
* github.com/microsoft/wmi: v0.34.0
* github.com/pkg/errors: v0.9.1
