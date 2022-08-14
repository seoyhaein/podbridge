package podbridge

import (
	"context"
	"errors"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/seoyhaein/utils"
)

// test

/*type (
	SetFunc  func(spec *specgen.SpecGenerator) *specgen.SpecGenerator
	SetFuncA func() *specgen.SpecGenerator
)*/

type Specgen *specgen.SpecGenerator

type CreateContainerResult struct {
	ErrorMessage error

	Name     string
	ID       string
	Warnings []string

	success bool
}

type ContainerSpec struct {
	spec *specgen.SpecGenerator
}

// TODO 컨테이너를 여러개 만들어야 하는 문제??
// TODO 꼼꼼히 테스트 해야함.
// pull goroutine

// TODO spec, pod 사용하지 않을 예정임. 대폭 수정예상

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

/*func NewSpec(imgName string) *specgen.SpecGenerator {
	if utils.IsEmptyString(imgName) {
		return nil
	}

	spec := specgen.NewSpecGenerator(imgName, false)

	return spec
}*/

//
func NewSpec() *ContainerSpec {
	return &ContainerSpec{
		spec: new(specgen.SpecGenerator),
	}
}

func (c *ContainerSpec) SetImage(imgName string) *ContainerSpec {
	if utils.IsEmptyString(imgName) {
		return nil
	}
	spec := specgen.NewSpecGenerator(imgName, false)
	c.spec = spec
	return c
}

func (c *ContainerSpec) SetOther(f func(spec Specgen) Specgen) *ContainerSpec {
	c.spec = f(c.spec)
	return c
}

// c.conSpec 가 파라미터로 들어가서 그 값을 세팅하는 func 를 외부에서 만든다.
/*func (c *ContainerSpec) SetOtherA(f func() *specgen.SpecGenerator) *ContainerSpec {
	// 데이터 병합이 문제다. 이걸 해결하자.
	// 포인터 문제가 발생한다. 젠장.

	//c.conSpec = f()
	//return c

	return func() *ContainerSpec {

		return nil
	}()
}*/

// TODO 수정해줘야 함.
func CreateContainer(ctx *context.Context, spec *specgen.SpecGenerator) *CreateContainerResult {
	var (
		result                 *CreateContainerResult
		containerExistsOptions containers.ExistsOptions
	)
	result = new(CreateContainerResult)
	err := spec.Validate()
	// TODO name, image 확인해야 할듯, 일 단 체크 해보자.
	if err != nil {
		result.ErrorMessage = err
		result.success = false
		return result
	}
	containerExistsOptions.External = utils.PFalse
	containerExists, err := containers.Exists(*ctx, spec.Name, &containerExistsOptions)
	if err != nil {
		result.ErrorMessage = err
		result.success = false
		return result
	}
	// 컨테이너가 local storage 에 존재하고 있다면
	if containerExists {
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = utils.PFalse
		containerData, err := containers.Inspect(*ctx, spec.Name, &containerInspectOptions)
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}
		if containerData.State.Running {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already running", spec.Name))
			result.ID = containerData.ID
			result.Name = spec.Name
			result.success = false
			return result
		} else {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already exists", spec.Name))
			result.ID = containerData.ID
			result.Name = spec.Name
			result.success = false
			return result
		}
	} else {
		imageExists, err := images.Exists(*ctx, spec.Image, nil)
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}
		// TODO 아래 코드는 필요 없을 듯, 이미지를 일단 만들어서 local 에 저장하는 구조임.
		if imageExists == false {
			_, err := images.Pull(*ctx, spec.Image, &images.PullOptions{})
			if err != nil {
				result.ErrorMessage = err
				result.success = false
				return result
			}
		}
		fmt.Printf("Pulling %s image...\n", Spec.Image)
		createResponse, err := containers.CreateWithSpec(*ctx, spec, &containers.CreateOptions{})
		if err != nil {
			result.ErrorMessage = err
			result.success = false
			return result
		}
		fmt.Printf("Creating %s container using %s image...\n", spec.Name, spec.Image)
		result.Name = spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings
	}
	result.success = true
	if Basket != nil {
		Basket.AddContainerId(result.ID)
	}
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

func (Res *CreateContainerResult) Kill(ctx *context.Context, options ...any) error {

	return nil
}

// 중요!! 지속적으로 container 의 status 를 확인해줘야 함으로, goroutine, 루프 구문이 들어가고 context 가 들어가고, channel 이 들어가야 할듯 하다.
// 퍼포먼스 문제가 있을까?? 일단 고민좀 해보자. 여러 컨테이너를 지속적으로 해야 함으로 이런 방식은 문제가 있을듯 하다. 일단 좀더 고민해 보자.

func (Res *CreateContainerResult) HealthCheck(ctx *context.Context, options ...any) error {

	//containers.RunHealthCheck()
	return nil
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
