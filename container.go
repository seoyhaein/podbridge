package podbridge

import (
	"context"
	"errors"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/seoyhaein/utils"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type SpecGen *specgen.SpecGenerator

type CreateContainerResult struct {
	ErrorMessage error

	Name     string
	ID       string
	Warnings []string

	// TODO 향후 int 로 바꿈.
	ContainerStatus string

	success bool
}

type ContainerSpec struct {
	Spec *specgen.SpecGenerator
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

// NewSpec
func NewSpec() *ContainerSpec {
	return &ContainerSpec{
		Spec: new(specgen.SpecGenerator),
	}
}

// SetImage
func (c *ContainerSpec) SetImage(imgName string) *ContainerSpec {
	if utils.IsEmptyString(imgName) {
		return nil
	}
	spec := specgen.NewSpecGenerator(imgName, false)
	c.Spec = spec
	return c
}

// SetOther
func (c *ContainerSpec) SetOther(f func(spec SpecGen) SpecGen) *ContainerSpec {
	c.Spec = f(c.Spec)
	return c
}

// CreateContainer
// TODO 수정해줘야 함.
func CreateContainer(ctx *context.Context, conSpec *ContainerSpec) *CreateContainerResult {
	var (
		result                 *CreateContainerResult
		containerExistsOptions containers.ExistsOptions
	)
	result = new(CreateContainerResult)
	err := conSpec.Spec.Validate()
	// TODO name, image 확인해야 할듯, 일 단 체크 해보자.
	if err != nil {
		//result.ErrorMessage = err
		//result.success = false
		panic(err)
		//return result
	}
	containerExistsOptions.External = utils.PFalse
	containerExists, err := containers.Exists(*ctx, conSpec.Spec.Name, &containerExistsOptions)
	if err != nil {
		//result.ErrorMessage = err
		//result.success = false
		//return result
		panic(err)
	}
	// 컨테이너가 local storage 에 존재하고 있다면
	if containerExists {
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = utils.PFalse
		containerData, err := containers.Inspect(*ctx, conSpec.Spec.Name, &containerInspectOptions)
		if err != nil {
			//result.ErrorMessage = err
			//result.success = false
			//return result
			panic(err)
		}
		if containerData.State.Running {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already running", conSpec.Spec.Name))
			result.ID = containerData.ID
			result.Name = conSpec.Spec.Name
			result.success = false
			result.ContainerStatus = "Running"
			return result
		} else {
			result.ErrorMessage = errors.New(fmt.Sprintf("%s container already exists", conSpec.Spec.Name))
			result.ID = containerData.ID
			result.Name = conSpec.Spec.Name
			result.success = false
			result.ContainerStatus = "Created"
			return result
		}
	} else {
		imageExists, err := images.Exists(*ctx, conSpec.Spec.Image, nil)
		if err != nil {
			panic(err)
			//result.ErrorMessage = err
			//result.success = false
			//return result
		}
		// TODO 아래 코드는 필요 없을 듯, 이미지를 일단 만들어서 local 에 저장하는 구조임.
		// basket 에 넣을지 고민하자.
		if imageExists == false {
			_, err := images.Pull(*ctx, conSpec.Spec.Image, &images.PullOptions{})
			if err != nil {
				panic(err)
				//result.ErrorMessage = err
				//result.success = false
				//return result
			}
		}
		Log.Infof("Pulling %s image...\n", conSpec.Spec.Image)
		createResponse, err := containers.CreateWithSpec(*ctx, conSpec.Spec, &containers.CreateOptions{})
		if err != nil {
			panic(err)
			//result.ErrorMessage = err
			//result.success = false
			//return result
		}
		Log.Infof("Creating %s container using %s image...\n", conSpec.Spec.Name, conSpec.Spec.Image)
		result.Name = conSpec.Spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings
	}
	result.success = true
	if Basket != nil {
		Basket.AddContainerId(result.ID)
	}
	return result
}

// Start
// startOptions 는 default 값을 사용한다.
// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerStartLibpod

func (Res *CreateContainerResult) Start(ctx *context.Context) error {
	if utils.IsEmptyString(Res.ID) == false && Res.ContainerStatus == "Created" {

		err := containers.Start(*ctx, Res.ID, &containers.StartOptions{})
		return err
	} else {
		return fmt.Errorf("cannot start container")
	}
}

// Stop
// TODO 추후 수정하자.
// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerStopLibpod
// default 값은 timeout 은  10 으로 세팅되어 있고, ignore 는 false 이다.
// ignore 는 만약 stop 된 컨테이너를 stop 되어 있을 때 stop 하는 경우 true 하면 에러 무시, false 로 하면 에러 리턴
// timeout 은 몇 후에 컨테어너를 kill 할지 정한다.
func (Res *CreateContainerResult) Stop(ctx *context.Context, options ...any) error {
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

// Kill
func (Res *CreateContainerResult) Kill(ctx *context.Context, options ...any) error {

	return nil
}

// HealthCheck
// 중요!! 지속적으로 container 의 status 를 확인해줘야 함으로, goroutine, 루프 구문이 들어가고 context 가 들어가고, channel 이 들어가야 할듯 하다.
// 퍼포먼스 문제가 있을까?? 일단 고민좀 해보자. 여러 컨테이너를 지속적으로 해야 함으로 이런 방식은 문제가 있을듯 하다. 일단 좀더 고민해 보자.
func (Res *CreateContainerResult) HealthCheck(ctx *context.Context, options ...any) error {

	//containers.RunHealthCheck()
	return nil
}

// 이미지 가존재하는지 확인하는 메서드 빼놓자.
// TODO wait 함수 구체적으로 살펴보기기
// 나머지들은 조금씩 구현해 나간다.
// containers.go
// TODO 중요 resource 관련
// https://github.com/containers/podman/issues/13145
// 명령어에 대한 heartbeat 관련 해서 처리 해야함.
// TODO 컨테이너의 상태를 확인하는 방법은 두가지 접근 방법이 있는데, local에 podman 이 설치 되어 있는 경우와, 원격(접속하는 머신에는 podman  이없음)에서 연결되는 경우
// 일단 먼저, local 에서 연결 하는 걸 적용한다. 구현하는 건 비교적 간단할 듯하다.
