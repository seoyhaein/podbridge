package podbridge

import (
	"reflect"

	nettypes "github.com/containers/common/libnetwork/types"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"
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

	pair struct {
		p1 any
		p2 any
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

// TODO 읽기 https://betterprogramming.pub/implementing-type-safe-tuples-with-go-1-18-9624010efaa
// https://github.com/golang/example/tree/master/gotypes
// 1.18 에서는 constraints: move to x/exp for Go 1.18 이렇게됨. 향후 조정될 수도 있음.
// https://github.com/golang/go/issues/50792
// 참고 : https://pkg.go.dev/go/types
// tuple 은 사용하지 않는다.

// TODO spec 테스트 하자.
// a 는 컨테이너나 pod 의 spec struct 이고, p 는 struct 에 넣을 필드와 값이다.
// 일단 string 만 적용되도록 했다.

func SetStringField(a interface{}, p pair) {

	specType := reflect.TypeOf(a)
	specValue := reflect.ValueOf(a)

	if specType.Kind() == reflect.Struct {
		//for _, p := range ps {
		v1, b1 := p.p1.(string)
		if b1 {
			_, find := specType.FieldByName(v1)

			if find {
				// 값을 설정 할 수 있다면...
				if specValue.FieldByName(v1).CanSet() {
					v2, b2 := p.p2.(string)
					if b2 {
						specValue.SetString(v2)
					}
				}
			}
		}
		//}
	}
}

// TODO test
func WithValues(a interface{}, ps ...pair) Option {
	return func(spec *specgen.SpecGenerator) Option {

		// for backup
		if spec == nil {
			s, b := a.(*specgen.SpecGenerator)

			if b {
				deepcopy.DeepCopy(Spec, s)
			}
		}

		backup = eraseSpec(backup)
		deepcopy.DeepCopy(backup, spec)

		for _, p := range ps {
			SetStringField(a, p)
		}

		return WithValues(backup)
	}
}
