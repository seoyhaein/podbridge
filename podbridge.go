package podbridge

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/containers/podman/v4/pkg/specgen"
)

/*
	전역적으로 재활용하면서 swap 할 2개의 포인터를 만들어 놓는다.
*/
var (
	Spec   *specgen.SpecGenerator
	backup *specgen.SpecGenerator
)

func init() {
	Spec = new(specgen.SpecGenerator)
	Spec.Name = "old hello world"
	Spec.Image = "docker.io/centos:latest"

	backup = new(specgen.SpecGenerator)
}

// TODO 에러에 관해서 좀 살펴보자.
// http://cloudrain21.com/golang-graceful-error-handling

type ResultCreateContainer struct {
	ErrorMessage string

	Name     string
	ID       string
	Warnings []string

	backup *specgen.SpecGenerator
}

// TODO 구현중
func CreatePod(ctx *context.Context, podId string) error {
	var podExistsOptions pods.ExistsOptions

	podExists, err := pods.Exists(*ctx, podId, &podExistsOptions)

	if err != nil {
		return err
	}
	// 기존에 pod 가 존재할 경우는 무조건 해당 pod 를 지운다.
	if podExists {
		return nil
	}

	return nil
}

// TODO 컨테이너를 여러개 만들어야 하는 문제??
// 리턴값을 통합하자. 지금은 일단 그냥 처리.
// 컨테이너를 생성하기만 한다.

func CreateContainerWithSpec(ctx *context.Context, options ...Option) *ResultCreateContainer {

	var (
		result *ResultCreateContainer
		old    Option
	)

	result = new(ResultCreateContainer)

	for _, opt := range options {
		if opt != nil {
			old = opt(Spec)
		}
	}

	var containerExistsOptions containers.ExistsOptions

	containerExistsOptions.External = PFalse
	containerExists, err := containers.Exists(*ctx, Spec.Name, &containerExistsOptions)

	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}

	// 컨테이너가 local storage 에 존재하고 있다면~
	if containerExists {
		// 참고, 다만 잘못된 정보일 수 있음.
		// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerInitLibpod
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = PFalse
		containerData, err := containers.Inspect(*ctx, Spec.Name, &containerInspectOptions)
		if err != nil {
			result.ErrorMessage = err.Error()
			return result
		}

		if containerData.State.Running {
			result.ErrorMessage = fmt.Sprintf("%s container already running", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			return result
		} else {
			result.ErrorMessage = fmt.Sprintf("%s container already exists", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			return result
		}
	} else {
		imageExists, err := images.Exists(*ctx, Spec.Name, nil)
		if err != nil {
			result.ErrorMessage = err.Error()
			return result
		}

		if !imageExists {
			_, err := images.Pull(*ctx, Spec.Image, nil)
			if err != nil {
				result.ErrorMessage = err.Error()
				return result
			}
		}

		fmt.Printf("Pulling %s image...\n", Spec.Image)

		createResponse, err := containers.CreateWithSpec(*ctx, Spec, nil)
		if err != nil {
			result.ErrorMessage = err.Error()
			return result
		}

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result.Name = Spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings
	}

	// default 값으로 저장
	old(nil)
	result.backup = Spec

	return result
}
