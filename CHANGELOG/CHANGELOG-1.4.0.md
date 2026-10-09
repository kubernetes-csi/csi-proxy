# Changelog since v1.3.0

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- Building csi-proxy and its client module now requires Go 1.26.

## Changes by Kind

### Feature

- None

### API

- None

### Bug or Regression

* Remove dependency on microsoft/wmi by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/415
* fix: skip GetSupportedSize when volume has not grown by @andyzhangx in https://github.com/kubernetes-csi/csi-proxy/pull/476

### Other (Cleanup or Flake)

* Release tool update by @darshansreenivas in https://github.com/kubernetes-csi/csi-proxy/pull/412
* fix: pin github action to exact SHA by @jsafrane in https://github.com/kubernetes-csi/csi-proxy/pull/414
* chore: enable dependabot for library-development branch by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/418
* Update SECURITY_CONTACTS with new contact information by @andyzhangx in https://github.com/kubernetes-csi/csi-proxy/pull/433
* chore: upgrade to go1.25 by @andyzhangx in https://github.com/kubernetes-csi/csi-proxy/pull/446
* Replace deprecated github.com/pkg/errors by @laozc in https://github.com/kubernetes-csi/csi-proxy/pull/459
* Update release-tools and gcb image to fix release builds by @torredil in https://github.com/kubernetes-csi/csi-proxy/pull/477

## Dependencies

### Added
_Nothing has changed._

### Changed

* github.com/google/go-cmp: v0.6.0 → v0.7.0
* github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus: v1.0.1 → v1.1.1
* github.com/prometheus/client_golang: v1.20.5 → v1.24.1
* github.com/sergi/go-diff: v1.3.1 → v1.4.0
* github.com/spf13/pflag: v1.0.5 → v1.0.10
* github.com/stretchr/testify: v1.10.0 → v1.12.1
* golang.org/x/sys: v0.32.0 → v0.48.0
* google.golang.org/grpc: v1.69.2 → v1.84.0 (fixes CVE-2026-84304)
* google.golang.org/protobuf: v1.36.0 → v1.36.12
* k8s.io/component-base: v0.28.4 → v0.37.1
* k8s.io/klog/v2: v2.130.1 → v2.140.0
* Bump actions/setup-go from 4 to 6 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/404
* Bump actions/checkout from 5 to 6 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/413
* Bump github.com/stretchr/testify from 1.10.0 to 1.11.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/419
* Bump k8s.io/klog/v2 from 2.130.1 to 2.140.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/421
* Bump google.golang.org/protobuf from 1.36.0 to 1.36.11 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/422
* Bump github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus from 1.0.1 to 1.1.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/423
* Bump github.com/spf13/pflag from 1.0.5 to 1.0.10 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/424
* Bump github.com/sergi/go-diff from 1.3.1 to 1.4.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/425
* Bump github.com/google/go-cmp from 0.6.0 to 0.7.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/427
* Bump actions/setup-go from 6.4.0 to 7.0.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/442
* Bump google.golang.org/grpc from 1.69.2 to 1.82.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/443
* Bump actions/checkout from 6.0.2 to 7.0.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/444
* Bump github.com/prometheus/client_golang from 1.20.5 to 1.24.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/445
* Bump github.com/prometheus/client_golang from 1.24.0 to 1.24.1 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/451
* Bump google.golang.org/grpc from 1.82.1 to 1.83.0 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/458
* Bump google.golang.org/protobuf from 1.36.11 to 1.36.12 by @dependabot[bot] in https://github.com/kubernetes-csi/csi-proxy/pull/462
* Upgrade dependencies in master by @torredil in https://github.com/kubernetes-csi/csi-proxy/pull/478

### Removed
* github.com/microsoft/wmi: v0.34.0
* github.com/pkg/errors: v0.9.1
