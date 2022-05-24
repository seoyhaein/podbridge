package podbridge

import (
	"context"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

func GetConnection() (*context.Context, error) {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")

	if sockDir == "" {
		sockDir = "/var/run"
	}

	socket := "unix:" + sockDir + "/podman/podman.sock"
	connText, err := bindings.NewConnection(context.Background(), socket)

	return &connText, err
}
