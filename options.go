package podbridge

import (
	nettypes "github.com/containers/common/libnetwork/types"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"
)

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
			}
		}

		backup = eraseSpec(backup)
		deepcopy.DeepCopy(backup, spec)

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

// TODO
// default 값으로 구체적인 값을 적용해놓자.
func InitBasicConfig() *BasicConfig {
	return &BasicConfig{
		Name: "new hello world",
	}
}

// TODO
// Spec 의 특정 필드만을 바꾸는 함수 필요
// 바꾸지만, 기본값은 본 상태로 돌려놓아야 한다.
func changeName(name string) *specgen.SpecGenerator {

	Spec.Name = name
	return Spec
}
