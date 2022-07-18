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

	// 이미지 만들기. 통합할 수 있는 함수 또는 메서드 만들자.
	store, err := pbr.NewBuildStore()

	if err != nil {
		return
	}

	builderOption := pbr.SetFromImage("alpine:latest")

	if builderOption == nil {
		return
	}

	builder, err := pbr.NewBuilder(*ctx, store, builderOption)

	imageId, err := pbr.BuildCustomImage(*ctx, builder, store, "localhost/helloWorld")

	if err != nil {
		return
	}

	fmt.Println("Image Id is : ", imageId)

	report, er = volumes.List(*ctx, &volumes.ListOptions{})

	if er != nil {
		for i, r := range report {
			fmt.Printf("%d: name:%s, driver:%s", i, r.Name, r.Driver)
		}
	}
}
