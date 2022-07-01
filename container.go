package podbridge

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
)

// TODO 에러에 관해서 좀 살펴보자.
// http://cloudrain21.com/golang-graceful-error-handling

type ResultCreateContainer struct {
	ErrorMessage string

	Name     string
	ID       string
	Warnings []string

	success bool
}

// TODO 컨테이너를 여러개 만들어야 하는 문제??
// TODO 중요 WithValues 를 새로 만들었기 때문에 CreateContainerWithSpec 수정 필요. 내일 하자.
// 컨테이너를 생성하기만 한다.
// 컨테이너 이름 자동생성

// 오류 찾기 NewSpecGenerator 과 비교해야 함

func StartContainerWithSpec(ctx *context.Context, conf *ContainerConfig) *ResultCreateContainer {

	if conf.IsSetSpec() == PFalse || conf.IsSetSpec() == nil {
		return nil
	}

	var (
		result                 *ResultCreateContainer
		containerExistsOptions containers.ExistsOptions
	)

	result = new(ResultCreateContainer)

	containerExistsOptions.External = PFalse
	containerExists, err := containers.Exists(*ctx, Spec.Name, &containerExistsOptions)

	if err != nil {
		result.ErrorMessage = err.Error()
		result.success = false
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
			result.success = false
			return result
		}

		if containerData.State.Running {
			result.ErrorMessage = fmt.Sprintf("%s container already running", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			result.success = false
			return result
		} else {
			result.ErrorMessage = fmt.Sprintf("%s container already exists", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			result.success = false
			return result
		}
	} else {
		// TODO 확인하자. 계속 pull 한다.
		imageExists, err := images.Exists(*ctx, Spec.Name, nil)
		if err != nil {
			result.ErrorMessage = err.Error()
			result.success = false
			return result
		}

		if !imageExists {
			_, err := images.Pull(*ctx, Spec.Image, nil)
			if err != nil {
				result.ErrorMessage = err.Error()
				result.success = false
				return result
			}
		}

		if conf.IsSetSpec() == PTrue {

			fmt.Printf("Pulling %s image...\n", Spec.Image)
			createResponse, err := containers.CreateWithSpec(*ctx, Spec, &containers.CreateOptions{})
			if err != nil {
				result.ErrorMessage = err.Error()
				result.success = false
				return result
			}

			fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

			result.Name = Spec.Name
			result.ID = createResponse.ID
			result.Warnings = createResponse.Warnings
		}

		// TODO 찾아보기  StartOptions
		err = containers.Start(*ctx, result.ID, &containers.StartOptions{})
		if err != nil {
			result.ErrorMessage = err.Error()
			result.success = false
			return result
		}
	}

	result.success = true
	return result
}
