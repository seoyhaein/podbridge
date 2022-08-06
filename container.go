package podbridge

import (
	"context"
	"errors"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/seoyhaein/utils"
)

type CreateContainerResult struct {
	ErrorMessage error

	Name     string
	ID       string
	Warnings []string

	success bool
}

// TODO 컨테이너를 여러개 만들어야 하는 문제??
// TODO 꼼꼼히 테스트 해야함.
// pull goroutine

func ContainerWithSpec(ctx *context.Context, conf *ContainerConfig) *CreateContainerResult {

	var (
		result                 *CreateContainerResult
		containerExistsOptions containers.ExistsOptions
	)

	result = new(CreateContainerResult)

	if conf.IsSetSpec() == utils.PFalse || conf.IsSetSpec() == nil {
		result.ErrorMessage = errors.New("Spec is not set")
		result.success = false
		return result
	}
	// 추가
	err := Spec.Validate()

	if err != nil {
		result.ErrorMessage = err
		result.success = false
		return result
	}

	containerExistsOptions.External = utils.PFalse
	containerExists, err := containers.Exists(*ctx, Spec.Name, &containerExistsOptions)

	if err != nil {
		result.ErrorMessage = err
		result.success = false
		return result
	}

	// 컨테이너가 local storage 에 존재하고 있다면
	if containerExists {
		// 참고, 다만 잘못된 정보일 수 있음.
		// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerInitLibpod
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = utils.PFalse
		containerData, err := containers.Inspect(*ctx, Spec.Name, &containerInspectOptions)
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}

		if containerData.State.Running {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already running", Spec.Name))
			result.ID = containerData.ID
			result.Name = Spec.Name
			result.success = false
			return result
		} else {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already exists", Spec.Name))
			result.ID = containerData.ID
			result.Name = Spec.Name
			result.success = false
			return result
		}
	} else {

		imageExists, err := images.Exists(*ctx, Spec.Image, nil)
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}

		// TODO 아래 코드는 필요 없을 듯, 이미지를 일단 만들어서 local 에 저장하는 구조임.
		if imageExists == false {
			_, err := images.Pull(*ctx, Spec.Image, &images.PullOptions{})
			if err != nil {
				result.ErrorMessage = err
				result.success = false
				return result
			}
		}

		fmt.Printf("Pulling %s image...\n", Spec.Image)

		createResponse, err := containers.CreateWithSpec(*ctx, Spec, &containers.CreateOptions{})
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result.Name = Spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings
	}

	result.success = true
	return result
}

func (Res *CreateContainerResult) Start(ctx *context.Context) error {
	// TODO 이 코드는 의미 없을것 같다. 테스트 할때 해보자.
	if Res == nil {
		return nil
	}

	if Res.success {
		// startOptions 는 default 값을 사용한다.
		// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerStartLibpod
		err := containers.Start(*ctx, Res.ID, &containers.StartOptions{})
		return err
	} else {
		return Res.ErrorMessage
	}
}

func (Res *CreateContainerResult) Stop(ctx *context.Context, options ...any) error {

	// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerStopLibpod
	// default 값은 timeout 은  10 으로 세팅되어 있고, ignore 는 false 이다.
	// ignore 는 만약 stop 된 컨테이너를 stop 되어 있을 때 stop 하는 경우 true 하면 에러 무시, false 로 하면 에러 리턴
	// timeout 은 몇 후에 컨테어너를 kill 할지 정한다.

	if Res == nil {
		return nil
	}
	stopOption := new(containers.StopOptions)
	for _, op := range options {
		v, b := op.(*bool)
		if b {
			stopOption.Ignore = v
		} else {
			v1, b1 := op.(*uint)
			if b1 {
				stopOption.Timeout = v1
			}
		}
	}

	err := containers.Stop(*ctx, Res.ID, stopOption)
	return err
}

// TODO wait 함수 구체적으로 살펴보기기
// 나머지들은 조금씩 구현해 나간다.
// containers.go

// TODO 중요 resource 관련
// https://github.com/containers/podman/issues/13145

// podbridge 에서 생성된 것만 지워야 한다.

// 명령어에 대한 heartbeat 관련 해서 처리 해야함.

// TODO 컨테이너의 상태를 확인하는 방법은 두가지 접근 방법이 있는데, local에 podman 이 설치 되어 있는 경우와, 원격(접속하는 머신에는 podman  이없음)에서 연결되는 경우
// 일단 먼저, local 에서 연결 하는 걸 적용한다. 구현하는 건 비교적 간단할 듯하다.
