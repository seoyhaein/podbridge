package podbridge

import (
	"context"
	"fmt"
	"testing"
)

// TODO 테스트 진행해야함. 오류 있음.

func TestStartContainerWithSpec(t *testing.T) {

	var (
		finally  Option
		finally1 Option
	)

	sockDir := DefaultLinuxSockDir()
	ctx, err := NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}

	conf := new(ContainerConfig)
	image := new(pair)

	image.p1 = "Image"
	image.p2 = "docker.io/busybox"

	opt := WithValues(image)
	finally = opt(Spec)

	isOk := conf.TrueAutoCreateContainerName(Spec)

	if isOk == nil {
		name := new(pair)
		name.p1 = "Name"
		name.p2 = "hello world"

		opt1 := WithValues(name)
		finally1 = opt1(Spec)
	}
	conf.TrueSetSpec()

	result := StartContainerWithSpec(ctx, conf)

	Finally(finally)

	if isOk == nil {
		Finally(finally1)
	}

	if result != nil {
		if result.success == false {
			fmt.Printf("error: %s", result.ErrorMessage)

		} else {
			fmt.Printf("Name: %s", result.Name)
			fmt.Printf("ID: %s", result.ID)
			fmt.Printf("Warnings: %s", result.Warnings)
		}

	}
}
