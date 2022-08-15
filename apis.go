package podbridge

import (
	"context"
)

// https://docs.docker.com/config/daemon/systemd/

// TODO 컨테이너 실행
func Run(ctx context.Context, imgeName string) (*CreateContainerResult, bool) {

	return nil, false
}

// 외부 노출되는 api 들을 여기에 정의해 놓자. 생각해보자.

//MustFirstCall used only in the init() function.
// mutex ListCreated 에 넣자.
// https://cloudolife.com/2020/04/18/Programming-Language/Golang-Go/Synchronization/Use-sync-Mutex-sync-RWMutex-to-lock-share-data-for-race-condition/
// Basket 은 singleton 이어야함 관련해서 처리해줘야 하고, 지금 은그냥 노출 하는데, api 를 통해서 노출하도록 처리한다.
// https://thebook.io/006806/ch05/03/03/
