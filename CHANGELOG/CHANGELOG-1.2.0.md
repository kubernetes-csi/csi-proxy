## Changes by Kind

### Feature

- Set min go version to 1.23 and godebug winsymlink=0 in go.mod ([#363](https://github.com/kubernetes-csi/csi-proxy/pull/363), [@mauriciopoppe](https://github.com/mauriciopoppe))
- The flag (-require-privacy) has been added. If true, New-SmbGlobalMapping will be called with -RequirePrivacy $true ([#315](https://github.com/kubernetes-csi/csi-proxy/pull/315), [@vitaliy-leschenko](https://github.com/vitaliy-leschenko))

### Bug or Regression

- Fix: unnecessary resize volume error ([#364](https://github.com/kubernetes-csi/csi-proxy/pull/364), [@andyzhangx](https://github.com/andyzhangx))
- Maximum retry (3) when for the getTarget method iterate and find the correct volume of the mount path. ([#336](https://github.com/kubernetes-csi/csi-proxy/pull/336), [@knabben](https://github.com/knabben))
- Total: 1 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 1, CRITICAL: 0)
  
  ┌──────────────────┬────────────────┬──────────┬────────┬───────────────────┬───────────────┬───────────────────────────────────────────────────┐
  │     Library      │ Vulnerability  │ Severity │ Status │ Installed Version │ Fixed Version │                       Title                       │
  ├──────────────────┼────────────────┼──────────┼────────┼───────────────────┼───────────────┼───────────────────────────────────────────────────┤
  │ golang.org/x/net │ CVE-2024-45338 │ HIGH     │ fixed  │ v0.32.0           │ 0.33.0        │ Non-linear parsing of case-insensitive content in │
  │                  │                │          │        │                   │               │ golang.org/x/net/html                             │
  │                  │                │          │        │                   │               │ https://avd.aquasec.com/nvd/cve-2024-45338        │
  └──────────────────┴────────────────┴──────────┴────────┴───────────────────┴───────────────┴───────────────────────────────────────────────────┘ ([#365](https://github.com/kubernetes-csi/csi-proxy/pull/365), [@andyzhangx](https://github.com/andyzhangx))

## Dependencies

### Added
- cel.dev/expr: v0.16.2
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp/tree/v1.24.2)
- github.com/go-logr/stdr: [v1.2.2](https://github.com/go-logr/stdr/tree/v1.2.2)
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0
- go.opentelemetry.io/otel/metric: v1.31.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0
- go.opentelemetry.io/otel/sdk: v1.31.0
- go.opentelemetry.io/otel/trace: v1.31.0
- go.opentelemetry.io/otel: v1.31.0
- golang.org/x/telemetry: bda5523

### Changed
- cloud.google.com/go/compute/metadata: v0.2.3 → v0.5.2
- github.com/Microsoft/go-winio: [v0.6.1 → v0.6.2](https://github.com/Microsoft/go-winio/compare/v0.6.1...v0.6.2)
- github.com/cespare/xxhash/v2: [v2.2.0 → v2.3.0](https://github.com/cespare/xxhash/v2/compare/v2.2.0...v2.3.0)
- github.com/cncf/udpa/go: [c52dc94 → 5459f2c](https://github.com/cncf/udpa/go/compare/c52dc94...5459f2c)
- github.com/cncf/xds/go: [e9ce688 → b4127c9](https://github.com/cncf/xds/go/compare/e9ce688...b4127c9)
- github.com/envoyproxy/go-control-plane: [9239064 → v0.13.1](https://github.com/envoyproxy/go-control-plane/compare/9239064...v0.13.1)
- github.com/envoyproxy/protoc-gen-validate: [v0.10.1 → v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v0.10.1...v1.1.0)
- github.com/go-logr/logr: [v1.2.0 → v1.4.2](https://github.com/go-logr/logr/compare/v1.2.0...v1.4.2)
- github.com/golang/glog: [v1.1.0 → v1.2.2](https://github.com/golang/glog/compare/v1.1.0...v1.2.2)
- github.com/golang/protobuf: [v1.5.3 → v1.5.4](https://github.com/golang/protobuf/compare/v1.5.3...v1.5.4)
- github.com/google/go-cmp: [v0.5.9 → v0.6.0](https://github.com/google/go-cmp/compare/v0.5.9...v0.6.0)
- github.com/google/uuid: [v1.3.0 → v1.6.0](https://github.com/google/uuid/compare/v1.3.0...v1.6.0)
- github.com/sirupsen/logrus: [v1.9.0 → v1.9.3](https://github.com/sirupsen/logrus/compare/v1.9.0...v1.9.3)
- github.com/stretchr/objx: [v0.5.0 → v0.5.2](https://github.com/stretchr/objx/compare/v0.5.0...v0.5.2)
- github.com/stretchr/testify: [v1.8.4 → v1.10.0](https://github.com/stretchr/testify/compare/v1.8.4...v1.10.0)
- golang.org/x/crypto: c2843e0 → v0.31.0
- golang.org/x/mod: v0.8.0 → v0.22.0
- golang.org/x/net: v0.9.0 → v0.33.0
- golang.org/x/oauth2: v0.7.0 → v0.23.0
- golang.org/x/sync: v0.1.0 → v0.10.0
- golang.org/x/sys: v0.11.0 → v0.28.0
- golang.org/x/term: v0.7.0 → v0.27.0
- golang.org/x/text: v0.9.0 → v0.21.0
- golang.org/x/tools: v0.6.0 → v0.28.0
- google.golang.org/appengine: v1.6.7 → v1.4.0
- google.golang.org/genproto/googleapis/api: dd9d682 → 796eee8
- google.golang.org/genproto/googleapis/rpc: 28d5490 → 9240e9c
- google.golang.org/genproto: 0005af6 → cb27e3a
- google.golang.org/grpc: v1.57.0 → v1.69.2
- google.golang.org/protobuf: v1.31.0 → v1.36.0
- k8s.io/klog/v2: v2.100.1 → v2.130.1

### Removed
- cloud.google.com/go/compute: v1.19.1
