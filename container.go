package podbridge

import (
	"context"
	"fmt"
	"time"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/specgen"
)

// TODO 에러에 관해서 좀 살펴보자.
// http://cloudrain21.com/golang-graceful-error-handling

type ResultCreateContainer struct {
	ErrorMessage string

	Name     string
	ID       string
	Warnings []string

	backup *specgen.SpecGenerator
}

// TODO 컨테이너를 여러개 만들어야 하는 문제??
// TODO 중요 WithValues 를 새로 만들었기 때문에 CreateContainerWithSpec 수정 필요. 내일 하자.
// 컨테이너를 생성하기만 한다.
// 컨테이너 이름 자동생성

func CreateContainerWithSpec(ctx *context.Context, options ...Option) (*ResultCreateContainer, bool) {

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
		return result, false
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
			return result, false
		}

		if containerData.State.Running {
			result.ErrorMessage = fmt.Sprintf("%s container already running", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			return result, false
		} else {
			result.ErrorMessage = fmt.Sprintf("%s container already exists", Spec.Name)
			result.ID = containerData.ID
			result.Name = Spec.Name
			return result, false
		}
	} else {
		imageExists, err := images.Exists(*ctx, Spec.Name, nil)
		if err != nil {
			result.ErrorMessage = err.Error()
			return result, false
		}

		if !imageExists {
			_, err := images.Pull(*ctx, Spec.Image, nil)
			if err != nil {
				result.ErrorMessage = err.Error()
				return result, false
			}
		}

		fmt.Printf("Pulling %s image...\n", Spec.Image)

		// 이름은 컨테이너를 생성할때 만들어준다.
		Spec.Name = createContainerName()

		createResponse, err := containers.CreateWithSpec(*ctx, Spec, nil)
		if err != nil {
			result.ErrorMessage = err.Error()
			return result, false
		}

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result.Name = Spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings

		// TODO 찾아보기  StartOptions
		err = containers.Start(*ctx, result.ID, &containers.StartOptions{})
		if err != nil {
			result.ErrorMessage = err.Error()
			return result, false
		}
	}

	// default 값으로 저장
	old(nil)
	result.backup = Spec

	return result, true
}

// TODO apis.go 로 이동 및 옵션을 만들어서 이름을 자동으로 만들어 줄지 설정할 수 있도록 한다.
// 일단 최초 컨테이너가 생성된 시점의 시간을 기록한다.
// 추가적으로 기록될 필요가 있는 정보가 있으면 추가한다.
func createContainerName() string {
	return time.Now().String()
}
