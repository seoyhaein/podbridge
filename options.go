package podbridge

import (
	"fmt"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v4/pkg/specgen"
)

// TODO 생각하기
// 이부분은 패키지를 달리 두는 방향으로 하고, 옵션을 여러개 둘수 있도록 한다.
// BasicConfig 와 관련해서 specgen.SpecGenerator 를 그대로 차용하는 건 어떨지 생각해보자.

type (
	Option func(*specgen.SpecGenerator) Option

	BasicConfig struct {
		Tag          string
		Name         string
		Image        string
		Volumes      []VolumeMount
		PortMappings []PortMapping
		Env          []EnvVar
		Command      []string

		// TODO 생각하자.
		Old *specgen.SpecGenerator
	}

	PortMapping struct {
		Text          string
		HostPort      uint16
		ContainerPort uint16
	}

	EnvVar struct {
		Text    string
		Key     string
		Value   string
		Mutable bool
	}

	VolumeMount struct {
		Text string
		Name string
		Dest string
	}
)

// 전역적으로 재활용하면서 swap 할 2개의 포인터를 만들어 놓는다.
var (
	Spec   *specgen.SpecGenerator
	Backup *specgen.SpecGenerator
)

func init() {
	Spec = new(specgen.SpecGenerator)
	Backup = new(specgen.SpecGenerator)
}

// deepcopy 문제 때문에 해당 struct 를 다 초기 세팅으로 만드는 수작업을 해야함.
func eraseSpec(spec *specgen.SpecGenerator) *specgen.SpecGenerator {
	return nil
}

// slice 관년 해서 참고하도록. append 하기 때문에 nil 로 해줘야 한다. https://yourbasic.org/golang/clear-slice/
//func WithBasic(basic *BasicConfig) Option {

// 입력 파라미터로 들어온 spec 에 혹시라도 값이 있으면 backup 한다.
// 포인터라 좀 버그 생각해야 한다. 포인터 복사할때 deepcopy 문제 발생할 개연성 존재. 테스트 해야함. 젠장.

// spec 은 사이즈가 큰 struct 인데, 컨테이너를 생성할때 반드시 필요한 struct 이다. 하지만, 컨테이너를 계속 생성하고 또한 지우고 하는 작업을 지속적으로 할때
// 향후 성능의 문제가 발생할 수 있을 것 같다. 따라서, 해당 spec 을 전역적으로 하나로 두고 이걸 재활용하는 방안을 생각해야한다.

// TODO 성능 테스트 반드시 필요. struct 재활용하는 것과 new 를 사용해서 재활용하지 않는 것.

func WithBasic(basic interface{}) Option {
	return func(spec *specgen.SpecGenerator) Option { // 여기 들어가는 spec 은 전역적으로 설정된 Spec 이다. 이건 항상 신규이다.
		oldbasic, isBasicConfig := basic.(*BasicConfig)

		// 버그 있다.
		if isBasicConfig {
			// oldbasic.Old 에 값이 존재하더라도 새로운 값으로 적용해준다.
			// 여기는 backup 용으로 제작된 field 이다.

			// TODO 여기서 부터 수정해야 함. 5/28
			// Spec, Backup
			// 먼저, Backup 을 Backup = eraseSpec(Backup) 해서 초기화 상태로 만들고,
			// Spec 또한 Spec = eraseSpec(Spec) 해준다. 이건 외부에서(options_test.go 5.28 todo 참고)
			// swap 하는 부분도 생각하자.
			// Backup = spec

			/*

				oldbasic.Old = new(specgen.SpecGenerator)
				oldbasic.Old = spec

				// 복사한 정보를 지운다. append 로 slice 를 채우기때문에 nil 처리한다.
				spec.PortMappings = nil
				spec.Env = nil
				spec.Volumes = nil

			*/

			spec.Name = oldbasic.Name
			for _, mapping := range oldbasic.PortMappings {
				spec.PortMappings = append(spec.PortMappings, nettypes.PortMapping{
					ContainerPort: mapping.ContainerPort,
					HostPort:      mapping.HostPort,
				})
			}

			if len(oldbasic.Command) > 0 {
				spec.Command = oldbasic.Command
			}

			if len(oldbasic.Env) > 0 {
				e := make(map[string]string)
				for _, env := range oldbasic.Env {
					e[env.Key] = env.Value
					spec.Env = e
				}
			}

			if len(oldbasic.Volumes) > 0 {
				for _, volume := range oldbasic.Volumes {
					vol := specgen.NamedVolume{
						Name: volume.Name,
						Dest: volume.Dest,
					}
					spec.Volumes = append(spec.Volumes, &vol)
				}
			}

			return WithBasic(oldbasic.Old)
		}

		oldspec, isSpec := basic.(*specgen.SpecGenerator)

		if isSpec {
			return WithBasic(oldspec)
		}

		return nil
	}
}

func WithTester(str string) Option {
	return func(spec *specgen.SpecGenerator) Option {
		fmt.Println(str)
		return WithTester("old")
	}
}
