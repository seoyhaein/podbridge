package podbridge

import (
	"context"
	"fmt"
	"testing"
)

// TODO 다양한 테스트 진행해야함.
func TestCreateContainerWithSpec(t *testing.T) {

	sockDir := DefaultLinuxSockDir()
	conText, err := NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}

	basicConfig := InitBasicConfig()
	opt := WithBasic(basicConfig)

	CreateContainerWithSpec(conText, opt)
}
