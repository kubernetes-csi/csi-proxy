#!/bin/bash

# This script uses docker buildx to build Windows container images
# on linux to run csi-proxy.exe inside HostProcess containers.
# Building Windows images on Linux requests setting --output=type=registry on the build commands.

set -o nounset
set -ex

: "${VERSION:?VERISON not set}"
: "${REGISTRY:?REGISTRY not set}"

export DOCKER_CLI_EXPERIMENTAL=enabled
if [ $(docker buildx ls | grep -c img-builder) == 0 ]; then
    docker buildx create --name img-builder 
fi
docker buildx use img-builder
docker buildx build --platform windows/amd64 --output=type=registry --pull -f Dockerfile --build-arg BASE=mcr.microsoft.com/windows/nanoserver:1809 -t $REGISTRY/csi-proxy:$VERSION-1809 ..
docker buildx build --platform windows/amd64 --output=type=registry --pull -f Dockerfile --build-arg BASE=mcr.microsoft.com/windows/nanoserver:ltsc2022 -t $REGISTRY/csi-proxy:$VERSION-ltsc2022 ..

docker manifest create $REGISTRY/csi-proxy:$VERSION $REGISTRY/csi-proxy:$VERSION-1809 $REGISTRY/csi-proxy:$VERSION-ltsc2022

os_version_1809=$(docker manifest inspect mcr.microsoft.com/windows/nanoserver:1809 | grep "os.version" | head -n 1 | awk -F\" '{print $4}')
docker manifest annotate --os windows --arch amd64 --os-version $os_version_1809 $REGISTRY/csi-proxy:$VERSION $REGISTRY/csi-proxy:$VERSION-1809

os_version_ltsc2022=$(docker manifest inspect mcr.microsoft.com/windows/nanoserver:ltsc2022 | grep "os.version" | head -n 1 | awk -F\" '{print $4}')
docker manifest annotate --os windows --arch amd64 --os-version $os_version_ltsc2022 $REGISTRY/csi-proxy:$VERSION $REGISTRY/csi-proxy:$VERSION-ltsc2022

docker manifest inspect $REGISTRY/csi-proxy:$VERSION
docker manifest push $REGISTRY/csi-proxy:$VERSION



