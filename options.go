package podbridge

import (
	nettypes "github.com/containers/common/libnetwork/types"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
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
		//Old *specgen.SpecGenerator
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

func eraseSpec(spec *specgen.SpecGenerator) *specgen.SpecGenerator {
	eraser := new(specgen.SpecGenerator)

	deepcopy.DeepCopy(spec, eraser)
	return spec
}

func WithBasic(basic interface{}) Option {
	return func(spec *specgen.SpecGenerator) Option {

		// for backup
		if spec == nil {
			s, b := basic.(*specgen.SpecGenerator)

			if b {
				deepcopy.DeepCopy(Spec, s)
				//Spec = s
			}
		}

		// 기존 spec 있던것을 Backup 에 넣는다.
		// Backup 을 초기화 상태로 만든다.
		backup = eraseSpec(backup)
		deepcopy.DeepCopy(backup, spec)
		//Backup = spec
		//spec = eraseSpec(spec)

		basicon, isBasicConfig := basic.(*BasicConfig)

		if isBasicConfig {

			spec.Name = basicon.Name
			for _, mapping := range basicon.PortMappings {
				spec.PortMappings = append(spec.PortMappings, nettypes.PortMapping{
					ContainerPort: mapping.ContainerPort,
					HostPort:      mapping.HostPort,
				})
			}

			if len(basicon.Command) > 0 {
				spec.Command = basicon.Command
			}

			if len(basicon.Env) > 0 {
				e := make(map[string]string)
				for _, env := range basicon.Env {
					e[env.Key] = env.Value
					spec.Env = e
				}
			}

			if len(basicon.Volumes) > 0 {
				for _, volume := range basicon.Volumes {
					vol := specgen.NamedVolume{
						Name: volume.Name,
						Dest: volume.Dest,
					}
					spec.Volumes = append(spec.Volumes, &vol)
				}
			}
		}

		return WithBasic(backup)
	}
}

// 테스트 용으로 제작
func InitBasicConfig() *BasicConfig {
	return &BasicConfig{
		Name: "new hello world",
	}
}

// defaultSpec 에서 생성된 Spec 의 특정 필드만을 바꾸는 함수 필요
func changeName(name string) *specgen.SpecGenerator {

	Spec.Name = name
	return Spec
}
