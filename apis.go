package podbridge

import (
	"context"
	"fmt"
)

// https://docs.docker.com/config/daemon/systemd/

// TODO 컨테이너 실행
func Run(ctx context.Context, imgeName string) (*CreateContainerResult, bool) {

	return nil, false
}

//MustFirstCall used only in the init() function.
// mutex ListCreated 에 넣자.
// https://cloudolife.com/2020/04/18/Programming-Language/Golang-Go/Synchronization/Use-sync-Mutex-sync-RWMutex-to-lock-share-data-for-race-condition/
// Basket 은 singleton 이어야함 관련해서 처리해줘야 하고, 지금 은그냥 노출 하는데, api 를 통해서 노출하도록 처리한다.
// https://thebook.io/006806/ch05/03/03/

func MustFirstCall() (*ListCreated, error) {
	basket, err := toListCreated()
	Basket = basket
	return Basket, err
}

func Save() error {
	if Basket == nil {
		fmt.Errorf("call MustFirstCall() first")
	}
	Basket.Save()
	return nil
}
