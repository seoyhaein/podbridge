package podbridge

import (
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

// 여기서 하나씩 추가할 수 있도록 한다.
// 테스트 진행하자.
// old 에러
// slice 관년 해서 참고하도록. append 하기 때문에 nil 로 해줘야 한다. https://yourbasic.org/golang/clear-slice/
// 테스트 진행하자.
// return 의 의미를 되짚어보자.
// TODO 파라미터가 BasicConfig 로 고정되는 문제점이 있다. 따라서 이부분도 생각해야 한다. interface, 생각해보기

func WithBasic(basic *BasicConfig) Option {
	return func(spec *specgen.SpecGenerator) Option {
		var old *BasicConfig
		// old 한 정보를 복사를 한다.
		// name
		old.Name = spec.Name
		//port
		for _, mapping := range spec.PortMappings {
			old.PortMappings = append(old.PortMappings, PortMapping{
				ContainerPort: mapping.ContainerPort,
				HostPort:      mapping.HostPort,
			})
		}
		// command
		old.Command = spec.Command
		// env
		// TODO old - env 일단 테스트 후 구현.
		// TODO old - volume 일단 테스트 후 구현.

		// 복사한 정보를 지운다. append 로 slice 를 채우기때문에 nil 처리한다.
		spec.PortMappings = nil
		spec.Env = nil
		spec.Volumes = nil

		spec.Name = basic.Name
		for _, mapping := range basic.PortMappings {
			spec.PortMappings = append(spec.PortMappings, nettypes.PortMapping{
				ContainerPort: mapping.ContainerPort,
				HostPort:      mapping.HostPort,
			})
		}

		if len(basic.Command) > 0 {
			spec.Command = basic.Command
		}

		if len(basic.Env) > 0 {
			e := make(map[string]string)
			for _, env := range basic.Env {
				e[env.Key] = env.Value
				spec.Env = e
			}
		}

		if len(basic.Volumes) > 0 {
			for _, volume := range basic.Volumes {
				vol := specgen.NamedVolume{
					Name: volume.Name,
					Dest: volume.Dest,
				}
				spec.Volumes = append(spec.Volumes, &vol)
			}
		}

		return WithBasic(old)
	}
}
