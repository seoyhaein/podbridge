#!/bin/sh

set -ex

git clone -b v4.1.0 https://github.com/containers/podman.git && cd podman
make binaries
cd bin/
copy podman /usr/local/bin
systemctl --user enable --now podman.socket && systemctl start --user podman.socket
podman system service -t 0 &