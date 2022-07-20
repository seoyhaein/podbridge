package podbridge

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/containers/buildah"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/storage/pkg/unshare"
)

func NewConnection(ctx context.Context, ipcName string) (*context.Context, error) {

	if len(strings.TrimSpace(ipcName)) == 0 {
		return nil, errors.New("ipcName cannot be an empty string")
	}
	conText, err := bindings.NewConnection(ctx, ipcName)

	return &conText, err
}

func defaultLinuxSockDir() (socket string) {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")
	if sockDir == "" {
		sockDir = "/var/run"
	}
	socket = "unix:" + sockDir + "/podman/podman.sock"

	return
}

func NewConnectionLinux(ctx context.Context, useBuildAh bool) (*context.Context, error) {
	socket := defaultLinuxSockDir()

	conText, err := bindings.NewConnection(ctx, socket)

	if useBuildAh {

		if buildah.InitReexec() {
			return nil, errors.New("InitReexec return false")
		}
		unshare.MaybeReexecUsingUserNamespace(false)
	}

	return &conText, err
}
