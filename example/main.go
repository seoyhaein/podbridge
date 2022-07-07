package main

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/volumes"
	"github.com/containers/podman/v4/pkg/domain/entities"
	pbr "github.com/seoyhaein/podbridge"
)

// 터미널에서 확인하기 위해서 만듬.
// go module 버그 있는 거 같다. 내가 잘못한 건가??

func main() {
	sockDir := pbr.DefaultLinuxSockDir()
	ctx, err := pbr.NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}
	var (
		report []*entities.VolumeListReport
		er     error
	)

	report, er = volumes.List(*ctx, &volumes.ListOptions{})

	if er == nil {
		for i, r := range report {
			fmt.Printf("%d: name:%s, driver:%s\n", i, r.Name, r.Driver)
		}
	} else {
		fmt.Println(er.Error())
	}
}
