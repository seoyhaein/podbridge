package main

import (
	"context"
	"fmt"

	pbr "github.com/seoyhaein/podbridge"
)

// 터미널에서 확인하기 위해서 만듬.

func main() {
	sockDir := pbr.DefaultLinuxSockDir()
	_, err := pbr.NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}

}
