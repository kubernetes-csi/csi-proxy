# Changelog since v1.0.2

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- None

## Changes by Kind

### Feature

- Support multiple working-dir parameters ([#184](https://github.com/kubernetes-csi/csi-proxy/pull/184), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Auto-start csi-proxy service on boot ([#194](https://github.com/kubernetes-csi/csi-proxy/pull/194), [@pradeep-hegde](https://github.com/pradeep-hegde))

### API
- v2alpha1 FileSystem API group with new API RmdirContents ([#186](https://github.com/kubernetes-csi/csi-proxy/pull/186), [@mauriciopoppe](https://github.com/mauriciopoppe))
- v2alpha1 Volume API group with new API GetClosestVolumeIDFromTargetPath ([#189](https://github.com/kubernetes-csi/csi-proxy/pull/189), [@mauriciopoppe](https://github.com/mauriciopoppe))

### Bug or Regression
- Reduce CSI proxy CPU usage ([#197](https://github.com/kubernetes-csi/csi-proxy/pull/197), [@pradeep-hegde](https://github.com/pradeep-hegde))

### Other (Cleanup or Flake)
- Update CSI Proxy dev scripts ([#182](https://github.com/kubernetes-csi/csi-proxy/pull/182), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Log working-dir parameters on startup ([#191](https://github.com/kubernetes-csi/csi-proxy/pull/191), [@mauriciopoppe](https://github.com/mauriciopoppe))
- Add codespell github action ([#198](https://github.com/kubernetes-csi/csi-proxy/pull/198), [@andyzhangx](https://github.com/andyzhangx))
