package podbridge

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/containers/podman/v4/pkg/bindings"
)

// return 값이 *context.Context 이라는 것을 주의하자.
// ipc 로 restful 연결 임으로 지속적으로 연결을 유지할 필요가 없지 않을까?
// 여기서는 client 입장에서 접근 한다.
// https://github.com/james-barrow/golang-ipc 참고
// 소스 정리전.

func GetConnection(ipcName string, ctx context.Context) (*context.Context, error) {
	// 주석처리된 코드는 필요없지만 일단은 주석 처리된 상태로 넣어 놓았다.
	/*
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
	*/
	/*sockDir := os.Getenv("XDG_RUNTIME_DIR")

	if sockDir == "" {
		sockDir = "/var/run"
	}

	socket := "unix:" + sockDir + "/podman/podman.sock"*/

	if len(strings.TrimSpace(ipcName)) == 0 {
		return nil, errors.New("ipcName cannot be an empty string")
	}
	connText, err := bindings.NewConnection(ctx, ipcName)

	return &connText, err
}

// help function
// podman 이 설치 되어 있는 것을 전제로 한다.
// 리눅스에서만...
func initSockDir() (socket string) {
	sockDir := os.Getenv("XDG_RUNTIME_DIR")
	if sockDir == "" {
		sockDir = "/var/run"
	}
	socket = "unix:" + sockDir + "/podman/podman.sock"

	return
}
