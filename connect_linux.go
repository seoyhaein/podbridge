package podbridge

import (
	"context"
	"log"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

func GetConnection() *context.Context {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")

	if sockDir == "" {
		sockDir = "/var/run"
	}
	socket := "unix:" + sockDir + "/podman/podman.sock"

	connText, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		log.Fatalln(err)
	}

	return &connText
}
