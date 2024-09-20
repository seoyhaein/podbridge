//This file now only builds on Linux.
//go:build linux
// +build linux

package podbridge

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/seoyhaein/utils"
)

func NewConnection5(ctx context.Context, ipcName string) (context.Context, error) {

	if utils.IsEmptyString(ipcName) {
		return nil, errors.New("ipcName cannot be an empty string")
	}
	cTx, err := bindings.NewConnection(ctx, ipcName)

	return cTx, err
}

func defaultLinuxSockDir5() (socket string) {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")
	if sockDir == "" {
		sockDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	socket = "unix:" + sockDir + "/podman/podman.sock"
	return
}

func NewConnectionLinux5(ctx context.Context) (context.Context, error) {

	socket := defaultLinuxSockDir5()
	cTx, err := bindings.NewConnection(ctx, socket)

	return cTx, err
}
