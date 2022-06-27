package podbridge

import (
	"fmt"

	nettypes "github.com/containers/common/libnetwork/types"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"

	"github.com/seoyhaein/go-tuple"
)

/*
	전역적으로 재활용하면서 swap 할 2개의 포인터를 만들어 놓는다.
*/
var (
	Spec   *specgen.SpecGenerator
	backup *specgen.SpecGenerator
)

// default 값을 세팅해놓았다.
func init() {
	Spec = new(specgen.SpecGenerator)
	Spec.Image = "alpine:latest"

	backup = new(specgen.SpecGenerator)
}

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

// TODO 중요! 정상작동하는지 테스트 필요.
func eraseSpec(spec *specgen.SpecGenerator) *specgen.SpecGenerator {
	eraser := new(specgen.SpecGenerator)

	deepcopy.DeepCopy(spec, eraser)
	return spec
}

func WithBasicConfig(basic interface{}) Option {
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

		return WithBasicConfig(backup)
	}
}

// TODO
// default 값으로 구체적인 값을 적용해놓자.
func InitBasicConfig() *BasicConfig {
	return &BasicConfig{
		Image: "docker.io/centos:latest",
	}
}

// tuple 적용
// TODO 읽기 https://betterprogramming.pub/implementing-type-safe-tuples-with-go-1-18-9624010efaa
// https://github.com/golang/example/tree/master/gotypes

func GenBasicConfig(as ...any) *BasicConfig {

	for _, a := range as {
		if a != nil {
			//if a
		}
	}

	return &BasicConfig{
		Image: "docker.io/centos:latest",
	}
}

// 1.18 에서는 constraints: move to x/exp for Go 1.18 이렇게됨. 향후 조정될 수도 있음.
// https://github.com/golang/go/issues/50792

func testuple() {

	// goland 버그 때문에 짜증난다.
	tup := tuple.New2(5, "hi!")
	fmt.Println(tup.V1) // Outputs 5.
	fmt.Println(tup.V2) // Outputs "hi!".
}
