package podbridge

import (
	"context"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

// return 값이 *context.Context 이라는 것을 주의하자.
// ipc 로 restful 연결 임으로 지속적으로 연결을 유지할 필요가 없지 않을까?

func GetConnection(ctx context.Context) (*context.Context, error) {
	// 주석처리된 코드는 필요없지만 일단은 주석 처리된 상태로 넣어 놓았다.
	/*
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
	*/
	sockDir := os.Getenv("XDG_RUNTIME_DIR")

	if sockDir == "" {
		sockDir = "/var/run"
	}

	socket := "unix:" + sockDir + "/podman/podman.sock"
	connText, err := bindings.NewConnection(ctx, socket)

	return &connText, err
}
