package podbridge

import (
	"context"
	"fmt"
	"testing"
)

// TODO 읽어보자
// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

func TestNewImages(t *testing.T) {

	ctx, err := NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	// 이미지 만들기. 통합할 수 있는 함수 또는 메서드 만들자.
	store, err := NewBuildStore()

	if err != nil {
		return
	}

	builderOption := SetFromImage("alpine:latest")

	if builderOption == nil {
		return
	}

	builder, err := NewBuilder(*ctx, store, builderOption)

	imageId, err := BuildCustomImage(*ctx, builder, store, "localhost/helloWorld")

	if err != nil {
		return
	}

	fmt.Println("Image Id is : ", imageId)

}
