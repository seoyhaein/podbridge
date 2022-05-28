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
	Option  func(*specgen.SpecGenerator) Option
	Recover func(*specgen.SpecGenerator) Option

	BasicConfig struct {
		Tag          string
		Name         string
		Image        string
		Volumes      []VolumeMount
		PortMappings []PortMapping
		Env          []EnvVar
		Command      []string

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

var recovery *specgen.SpecGenerator

func init() {
	recovery = new(specgen.SpecGenerator)
}

// 여기서 하나씩 추가할 수 있도록 한다.
// 테스트 진행하자.
// old 에러
// slice 관년 해서 참고하도록. append 하기 때문에 nil 로 해줘야 한다. https://yourbasic.org/golang/clear-slice/
// 테스트 진행하자.
// return 의 의미를 되짚어보자.

// 포인터로 될까????

//func WithBasic(basic *BasicConfig) Option {
func WithBasic(basic interface{}) Option {
	return func(spec *specgen.SpecGenerator) Option {
		oldbasic, isBasicConfig := basic.(*BasicConfig)

		// true 이면
		if isBasicConfig {
			oldbasic.Old = new(specgen.SpecGenerator)
			oldbasic.Old = spec

			// 복사한 정보를 지운다. append 로 slice 를 채우기때문에 nil 처리한다.
			spec.PortMappings = nil
			spec.Env = nil
			spec.Volumes = nil

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
