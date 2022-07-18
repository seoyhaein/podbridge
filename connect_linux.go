package podbridge

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/containers/podman/v4/pkg/bindings"
)

func NewConnection(ctx context.Context, ipcName string) (*context.Context, error) {

	if len(strings.TrimSpace(ipcName)) == 0 {
		return nil, errors.New("ipcName cannot be an empty string")
	}
	connText, err := bindings.NewConnection(ctx, ipcName)

	return &connText, err
}

func DefaultLinuxSockDir() (socket string) {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")
	if sockDir == "" {
		sockDir = "/var/run"
	}
	socket = "unix:" + sockDir + "/podman/podman.sock"

	return
}

func NewConnectionLinux(ctx context.Context) (*context.Context, error) {
	socket := DefaultLinuxSockDir()

	ctx, err := bindings.NewConnection(ctx, socket)
	return &ctx, err
}
