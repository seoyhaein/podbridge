package podbridge

import (
	"context"
	"fmt"
	"time"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/seoyhaein/utils"
	"github.com/sirupsen/logrus"
)

// TODO *context.Context -> context.Context 로 변경한다.

var Log = logrus.New()

type (
	SpecGen *specgen.SpecGenerator

	CreateContainerResult struct {
		Name     string
		ID       string
		Warnings []string
		Status   ContainerStatus
		ch       chan ContainerStatus
	}

	ContainerSpec struct {
		Spec *specgen.SpecGenerator
	}
)

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

// SetHealthChecker
func (c *ContainerSpec) SetHealthChecker(inCmd, interval string, retries uint, timeout, startPeriod string) *ContainerSpec {
	// cf. SetHealthChecker("CMD-SHELL /app/healthcheck.sh", "2s", 3, "30s", "1s")
	healthConfig, err := SetHealthChecker(inCmd, interval, retries, timeout, startPeriod)
	if err != nil {
		panic(err)
	}
	c.Spec.HealthConfig = healthConfig
	return c
}

// CreateContainer
// TODO 수정해줘야 함. name, image 확인해야 할듯, 일단 체크 해보자.
func CreateContainer(ctx *context.Context, conSpec *ContainerSpec) *CreateContainerResult {
	var (
		result                 *CreateContainerResult
		containerExistsOptions containers.ExistsOptions
	)
	result = new(CreateContainerResult)
	err := conSpec.Spec.Validate()
	if err != nil {
		panic(err)
	}
	containerExistsOptions.External = utils.PFalse
	containerExists, err := containers.Exists(*ctx, conSpec.Spec.Name, &containerExistsOptions)
	if err != nil {
		panic(err)
	}
	// 컨테이너가 local storage 에 존재하고 있다면
	if containerExists {
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = utils.PFalse
		containerData, err := containers.Inspect(*ctx, conSpec.Spec.Name, &containerInspectOptions)
		if err != nil {
			panic(err)
		}
		if containerData.State.Running {
			Log.Infof("%s container already running", conSpec.Spec.Name)
			result.ID = containerData.ID
			result.Name = conSpec.Spec.Name
			result.Status = Running
			return result
		} else {
			Log.Infof("%s container already exists", conSpec.Spec.Name)
			result.ID = containerData.ID
			result.Name = conSpec.Spec.Name
			result.Status = Created
			return result
		}
	} else {
		imageExists, err := images.Exists(*ctx, conSpec.Spec.Image, nil)
		if err != nil {
			panic(err)
		}
		// TODO basket 에 넣을지 고민하자.
		if imageExists == false {
			_, err := images.Pull(*ctx, conSpec.Spec.Image, &images.PullOptions{})
			if err != nil {
				panic(err)
			}
		}
		Log.Infof("Pulling %s image...\n", conSpec.Spec.Image)
		createResponse, err := containers.CreateWithSpec(*ctx, conSpec.Spec, &containers.CreateOptions{})
		if err != nil {
			panic(err)
		}
		Log.Infof("Creating %s container using %s image...\n", conSpec.Spec.Name, conSpec.Spec.Image)
		result.Name = conSpec.Spec.Name
		result.ID = createResponse.ID
		result.Warnings = createResponse.Warnings
		result.Status = Created
	}
	if Basket != nil {
		Basket.AddContainerId(result.ID)
	}
	return result
}

// Start
// startOptions 는 default 값을 사용한다.
// https://docs.podman.io/en/latest/_static/api.html?version=v4.1#operation/ContainerStartLibpod
func (Res *CreateContainerResult) Start(ctx *context.Context) error {
	if utils.IsEmptyString(Res.ID) == false && Res.Status == Created {
		err := containers.Start(*ctx, Res.ID, &containers.StartOptions{})
		return err
	} else {
		return fmt.Errorf("cannot start container")
	}
}

// ReStart 중복되는 것 같긴하다. 수정해줘야 한다. ReStart
func (Res *CreateContainerResult) ReStart(ctx *context.Context) error {
	if utils.IsEmptyString(Res.ID) == false && Res.Status != Running {
		err := containers.Start(*ctx, Res.ID, &containers.StartOptions{})
		return err
	} else {
		return fmt.Errorf("cannot re-start container")
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
			v1, b1 := op.(uint)
			if b1 {
				stopOption.Timeout = &v1
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

// https://github.com/containers/podman/issues/12226
// https://developers.redhat.com/blog/2019/04/18/monitoring-container-vitality-and-availability-with-podman#interacting_with_the_results_of_healthchecks
// https://devops.stackexchange.com/questions/11501/healthcheck-cmd-vs-cmd-shell
// https://nomad-programmer.tistory.com/309
// https://www.codegrepper.com/code-examples/shell/how+to+check+if+a+process+is+running+in+linux+using+shell+script
// https://www.redhat.com/sysadmin/error-handling-bash-scripting
// https://twpower.github.io/134-how-to-return-shell-scipt-value

// healthchecker shell 의 경우는 환경변수를 만들고, 여기에, shell script 진행 상황에 따라 환경변수를 집어넣는 방식으로 진행한다.
// TODO 테스트는 성공했는데, 보강해야할 것들이 많다. healthy 및 기타 status 에 따라 종료 되는 것과, 컨테이너 자체의 종료를 알아야 한다.
func (Res *CreateContainerResult) HealthCheck(ctx *context.Context, interval string) error {
	// TODO 잘 살펴보자
	// sender, close ???, cancel 테스트 하자.
	go func(ctx context.Context, res *CreateContainerResult) {
		if res.ch == nil {
			// TODO 일단 buffer 를 100 으로 주었다.
			res.ch = make(chan ContainerStatus, 100)
		}
		intervalDuration, err := time.ParseDuration(interval)
		if err != nil {
			intervalDuration = time.Second
		}
		ticker := time.Tick(intervalDuration)
		for {
			select {
			case <-ticker:
				healthCheck, err := containers.RunHealthCheck(ctx, res.ID, &containers.HealthCheckOptions{})
				if err != nil {
					break
				}
				if healthCheck.Status == "healthy" {
					res.ch <- Healthy
				}

				if healthCheck.Status == "unhealthy" {
					res.ch <- Unhealthy
				}
			case <-ctx.Done():
				close(res.ch)
				break
			}
		}
	}(*ctx, Res)

	// sender, receiver 를 여기서 구현...
	// sender 의 경우는 goroutine 으로
	// receiver 의 경우는 그냥, context 처리, 여기서 기다려준다.

	return nil
}

// Run CreateContainer, Start or Restart, HealthCheck 들어가고 Receiver 들어가는 메서드
// unhealthy 이거나 container 가 중단(정상종료, 비정상 종료)될 경우 멈춤.
// ctx.Err()
// Exit 코드도 잡아야 함. 그럼 inspect 로 처리하는게 나을듯.
func (Res *CreateContainerResult) Run(ctx context.Context) (*CreateContainerResult, error) {
	if Res.ch == nil {
		return nil, fmt.Errorf("channel not created")
	}

	for {
		select {
		case status, ok := <-Res.ch:
			if ok {
				if status == Unhealthy {
					return nil, nil
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
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
