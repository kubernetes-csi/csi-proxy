# Changelog since v1.1.2

## Bug or Regression

- Ensuring all paths in powershell are passes as env vars ([#306](https://github.com/kubernetes-csi/csi-proxy/pull/306), [@ritazh](https://github.com/ritazh))

## Dependencies

### Added
- cloud.google.com/go/compute/metadata: v0.2.3
- cloud.google.com/go/compute: v1.19.1
- github.com/cespare/xxhash/v2: [v2.2.0](https://github.com/cespare/xxhash/v2/tree/v2.2.0)
- github.com/cncf/xds/go: [e9ce688](https://github.com/cncf/xds/go/tree/e9ce688)
- github.com/kr/pretty: [v0.1.0](https://github.com/kr/pretty/tree/v0.1.0)
- github.com/kr/pty: [v1.1.1](https://github.com/kr/pty/tree/v1.1.1)
- github.com/kr/text: [v0.1.0](https://github.com/kr/text/tree/v0.1.0)
- github.com/yuin/goldmark: [v1.4.13](https://github.com/yuin/goldmark/tree/v1.4.13)
- golang.org/x/mod: v0.8.0
- golang.org/x/term: v0.7.0
- google.golang.org/genproto/googleapis/api: dd9d682
- google.golang.org/genproto/googleapis/rpc: 28d5490
- gopkg.in/yaml.v3: v3.0.1

### Changed
- github.com/Microsoft/go-winio: [v0.4.16 → v0.6.1](https://github.com/Microsoft/go-winio/compare/v0.4.16...v0.6.1)
- github.com/census-instrumentation/opencensus-proto: [v0.2.1 → v0.4.1](https://github.com/census-instrumentation/opencensus-proto/compare/v0.2.1...v0.4.1)
- github.com/cncf/udpa/go: [5459f2c → c52dc94](https://github.com/cncf/udpa/go/compare/5459f2c...c52dc94)
- github.com/envoyproxy/go-control-plane: [668b12f → 9239064](https://github.com/envoyproxy/go-control-plane/compare/668b12f...9239064)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0 → v0.10.1](https://github.com/envoyproxy/protoc-gen-validate/compare/v0.1.0...v0.10.1)
- github.com/go-logr/logr: [v0.4.0 → v1.2.0](https://github.com/go-logr/logr/compare/v0.4.0...v1.2.0)
- github.com/golang/glog: [23def4e → v1.1.0](https://github.com/golang/glog/compare/23def4e...v1.1.0)
- github.com/golang/protobuf: [v1.4.3 → v1.5.3](https://github.com/golang/protobuf/compare/v1.4.3...v1.5.3)
- github.com/google/go-cmp: [v0.5.0 → v0.5.9](https://github.com/google/go-cmp/compare/v0.5.0...v0.5.9)
- github.com/google/uuid: [v1.1.2 → v1.3.0](https://github.com/google/uuid/compare/v1.1.2...v1.3.0)
- github.com/iancoleman/strcase: [e506e3e → v0.3.0](https://github.com/iancoleman/strcase/compare/e506e3e...v0.3.0)
- github.com/sergi/go-diff: [v1.0.0 → v1.3.1](https://github.com/sergi/go-diff/compare/v1.0.0...v1.3.1)
- github.com/sirupsen/logrus: [v1.4.1 → v1.9.0](https://github.com/sirupsen/logrus/compare/v1.4.1...v1.9.0)
- github.com/stretchr/objx: [v0.1.1 → v0.5.0](https://github.com/stretchr/objx/compare/v0.1.1...v0.5.0)
- github.com/stretchr/testify: [v1.5.1 → v1.8.4](https://github.com/stretchr/testify/compare/v1.5.1...v1.8.4)
- golang.org/x/net: 1617124 → v0.9.0
- golang.org/x/oauth2: d2e6202 → v0.7.0
- golang.org/x/sync: 1122301 → v0.1.0
- golang.org/x/sys: d101bd2 → v0.11.0
- golang.org/x/text: v0.3.2 → v0.9.0
- golang.org/x/tools: 2c0ae70 → v0.6.0
- google.golang.org/appengine: v1.4.0 → v1.6.7
- google.golang.org/genproto: cb27e3a → 0005af6
- google.golang.org/grpc: v1.38.0 → v1.57.0
- google.golang.org/protobuf: v1.25.0 → v1.31.0
- gopkg.in/check.v1: 20d25e2 → 41f04d3
- gopkg.in/yaml.v2: v2.2.2 → v2.4.0
- k8s.io/klog/v2: v2.9.0 → v2.100.1

### Removed
_Nothing has changed._
