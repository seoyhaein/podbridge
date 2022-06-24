package podbridge

import (
	"fmt"

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

/*
	전역적으로 재활용하면서 swap 할 2개의 포인터를 만들어 놓는다.
*/
var (
	Spec   *specgen.SpecGenerator
	Backup *specgen.SpecGenerator
)

func init() {
	Spec = new(specgen.SpecGenerator)
	Spec = defaultSpec(Spec)

	Backup = new(specgen.SpecGenerator)
}

// deepcopy 문제 때문에 해당 struct 를 다 초기 세팅으로 만드는 수작업을 해야함.
// 수작업으로 꼭해야 하나???
func eraseSpec(spec *specgen.SpecGenerator) *specgen.SpecGenerator {
	return spec
}

//func WithBasic(basic *BasicConfig) Option {

// TODO 성능 테스트 반드시 필요. struct 재활용하는 것과 new 를 사용해서 재활용하지 않는 것.
// spec 은 사이즈가 큰 struct 인데, 컨테이너를 생성할때 반드시 필요한 struct 이다. 하지만, 컨테이너를 계속 생성하고 또한 지우고 하는 작업을 지속적으로 할때
// 향후 성능의 문제가 발생할 수 있을 것 같다. 따라서, 해당 spec 을 전역적으로 하나로 두고 이걸 재활용하는 방안을 생각해야한다.
// basic 은 신규로 들어가는 녀석이고, spec 은 default 값이다.
// 신규로 spec 을 만들면 안된다.
// spec 에 들어가는 녀석은 Spec 이다.
// 초기 들어가는 Spec 은 default 이다.

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
		//Backup = eraseSpec(Backup)
		// TODO deepcopy 문제 살펴보기.
		deepcopy.DeepCopy(Backup, spec)
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

		return WithBasic(Backup)
	}
}

func WithTester(str string) Option {
	return func(spec *specgen.SpecGenerator) Option {
		fmt.Println(str)
		return WithTester("old")
	}
}

// 테스트 용으로 제작
func InitBasicConfig() *BasicConfig {
	return &BasicConfig{
		Name:  "new hello world",
		Image: "docker.io/ubuntu:latest",
	}
}

// 처음사용되고 끝
func defaultSpec(spec *specgen.SpecGenerator) *specgen.SpecGenerator {

	if spec == nil {
		return nil
	}

	spec.Name = "old hello world"
	spec.Image = "docker.io/centos:latest"

	return spec
}

// defaultSpec 에서 생성된 Spec 의 특정 필드만을 바꾸는 함수 필요
func changeName(name string) *specgen.SpecGenerator {

	Spec.Name = name
	return Spec
}
