package main

import (
	"context"
	"fmt"

	v4en "github.com/containers/podman/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/bindings/volumes"
	pbr "github.com/seoyhaein/podbridge"
)

// 터미널에서 확인하기 위해서 만듬.

func main() {
	sockDir := pbr.DefaultLinuxSockDir()
	ctx, err := pbr.NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}
	var (
		report []*v4en.VolumeListReport
		er     error
	)

	report, er = volumes.List(*ctx, &volumes.ListOptions{})

	if er != nil {
		for i, r := range report {
			fmt.Sprintf("%d: name:%s, driver:%s", i, r.Name, r.Driver)
		}
	}
}
