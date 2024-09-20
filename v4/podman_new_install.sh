#!/bin/bash

set -ex

# 최신 버전 변수
LATEST_PODMAN_VERSION="4.1.0"
LATEST_RUNC_VERSION="1.1.3"
LATEST_CONMON_VERSION="2.1.3"

# versions.env 파일이 있는지 확인하고 로드
if [ -f "versions.env" ]; then
    source versions.env
else
    echo "versions.env 파일을 찾을 수 없습니다. 기본 버전을 사용합니다."
    PODMAN_VERSION=$LATEST_PODMAN_VERSION
    RUNC_VERSION=$LATEST_RUNC_VERSION
    CONMON_VERSION=$LATEST_CONMON_VERSION
fi

# Podman이 설치되어 있는지 확인
INSTALL_PODMAN=false
if command -v podman &> /dev/null; then
    # 설치된 Podman 버전 확인
    INSTALLED_PODMAN_VERSION=$(podman version --format '{{.Version}}')
    echo "Installed Podman version: $INSTALLED_PODMAN_VERSION"

    # 설치된 버전이 최신 버전인지 확인
    if [ "$INSTALLED_PODMAN_VERSION" != "$PODMAN_VERSION" ]; then
        echo "Updating Podman from version $INSTALLED_PODMAN_VERSION to $PODMAN_VERSION..."
        # podman이 실행 중인 경우 종료
        echo "Stopping running Podman processes..."
        podman ps -q | xargs -r podman stop
        # Podman API 서버 중지
        sudo systemctl stop podman.service
        sudo systemctl stop podman.socket

        # Podman API 서버 중지 및 비활성화
        sudo systemctl disable --now podman.socket
        sudo systemctl disable --now podman.service

        sleep 2  # 잠시 대기하여 서비스가 완전히 중지되도록 함
        INSTALL_PODMAN=true
    else
        echo "Podman is already up-to-date."
    fi
else
    echo "Podman is not installed. Proceeding with installation..."
    INSTALL_PODMAN=true
fi

# runc가 설치되어 있는지 확인
INSTALL_RUNC=false
if command -v runc &> /dev/null; then
    # 설치된 runc 버전 확인
    INSTALLED_RUNC_VERSION=$(runc --version | grep runc | awk '{print $3}')
    echo "Installed runc version: $INSTALLED_RUNC_VERSION"

    # 설치된 버전이 최신 버전인지 확인
    if [ "$INSTALLED_RUNC_VERSION" != "$RUNC_VERSION" ]; then
        echo "Updating runc from version $INSTALLED_RUNC_VERSION to $RUNC_VERSION..."
        INSTALL_RUNC=true
    else
        echo "runc is already up-to-date."
    fi
else
    echo "runc is not installed. Proceeding with installation..."
    INSTALL_RUNC=true
fi

# conmon이 설치되어 있는지 확인
INSTALL_CONMON=false
if command -v conmon &> /dev/null; then
    # 설치된 conmon 버전 확인
    INSTALLED_CONMON_VERSION=$(conmon --version | grep conmon | awk '{print $3}')
    echo "Installed conmon version: $INSTALLED_CONMON_VERSION"

    # 설치된 버전이 최신 버전인지 확인
    if [ "$INSTALLED_CONMON_VERSION" != "$CONMON_VERSION" ]; then
        echo "Updating conmon from version $INSTALLED_CONMON_VERSION to $CONMON_VERSION..."
        INSTALL_CONMON=true
    else
        echo "conmon is already up-to-date."
    fi
else
    echo "conmon is not installed. Proceeding with installation..."
    INSTALL_CONMON=true
fi

# 의존성 설치
sudo apt-get update
sudo apt-get install -y \
  gcc \
  btrfs-progs \
  go-md2man \
  iptables \
  libassuan-dev \
  libbtrfs-dev \
  libc6-dev \
  libdevmapper-dev \
  libglib2.0-dev \
  libgpgme-dev \
  libgpg-error-dev \
  libprotobuf-dev \
  libprotobuf-c-dev \
  libseccomp-dev \
  libselinux1-dev \
  libsystemd-dev \
  pkg-config \
  uidmap

# runc 설치
if [ "$INSTALL_RUNC" = true ]; then
    pushd ~/go/src/github.com/seoyhaein/
    rm -rf runc
    git clone -b v$RUNC_VERSION https://github.com/opencontainers/runc.git && pushd runc
    make
    sudo cp runc /usr/local/bin
    popd
    popd
fi

# conmon 설치
if [ "$INSTALL_CONMON" = true ]; then
    pushd ~/go/src/github.com/seoyhaein/
    rm -rf conmon
    git clone -b v$CONMON_VERSION https://github.com/containers/conmon.git && pushd conmon
    make
    pushd bin/
    sudo cp conmon /usr/local/bin
    popd
    popd
fi

# podman 설치
if [ "$INSTALL_PODMAN" = true ]; then
    # Podman 서비스 중지
    sudo systemctl stop podman.socket || true
    sudo systemctl stop podman.service || true
    pushd ~/go/src/github.com/seoyhaein/
    rm -rf podman
    git clone -b v$PODMAN_VERSION https://github.com/containers/podman.git && pushd podman
    make binaries
    pushd bin/
    sudo cp podman /usr/local/bin
    popd
    popd

    # Podman 서비스 다시 시작
    sudo systemctl start podman.socket
    sudo systemctl start podman.service
    sudo systemctl --user enable --now podman.socket && sudo systemctl start --user podman.socket
    podman system service -t 0 &
fi

cd ~/go/src/github.com/seoyhaein/podbridge
