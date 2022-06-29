package podbridge

import (
	"errors"
	"fmt"
	"reflect"

	nettypes "github.com/containers/common/libnetwork/types"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"
)

// 전역적으로 재활용하면서 swap 할 2개의 포인터를 만들어 놓는다.
var (
	Spec   *specgen.SpecGenerator
	backup *specgen.SpecGenerator
)

// default 값을 세팅해놓았다.
func init() {
	Spec = new(specgen.SpecGenerator)
	//TODO 삭제할지 고민 중
	Spec.Image = "alpine:latest"

	backup = new(specgen.SpecGenerator)
}

type (
	Option func(*specgen.SpecGenerator) Option

	// deprecated
	BasicConfig struct {
		Tag          string
		Name         string
		Image        string
		Volumes      []VolumeMount
		PortMappings []PortMapping
		Env          []EnvVar
		Command      []string
	}
	// deprecated
	PortMapping struct {
		Text          string
		HostPort      uint16
		ContainerPort uint16
	}
	// deprecated
	EnvVar struct {
		Text    string
		Key     string
		Value   string
		Mutable bool
	}
	// deprecated
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
	clear := clearStruct(spec)

	if clear == nil {
		return nil
	}
	s, b := clear.(*specgen.SpecGenerator)

	if b {
		return s
	} else {
		return nil
	}

}

/*
	If a reflect.Value is a pointer, then v.Elem() is equivalent to reflect.Indirect(v). If it is not a pointer, then they are not equivalent:
		If the value is an interface then reflect.Indirect(v) will return the same value, while v.Elem() will return the contained dynamic value.
		If the value is something else, then v.Elem() will panic.
	The reflect.Indirect helper is intended for cases where you want to accept either a particular type, or a pointer to that type.
	One example is the database/sql conversion routines: by using reflect.Indirect, it can use the same code paths to handle the various types and pointers to those types.
*/

func clearStruct(a interface{}) interface{} {

	v := reflect.Indirect(reflect.ValueOf(a))
	if v.Kind() == reflect.Struct {
		v.Set(reflect.Zero(v.Type()))

		return v
	}

	return nil
}

// deprecated
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

// deprecated
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
// a 는 컨테이너나 pod 의 spec struct 이고, p 는 struct 에 넣을 필드와 값이다.
// 일단 string 만 적용되도록 했다.

// deprecated
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
// TODO 이런 함수의 형태는 error 처리가 좀 힘든데, 어떻할지 생각해보자. 일단 panic 으로 설정함.

func WithValues(a interface{}, ps ...*pair) Option {
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
			v1, b1 := p.p1.(string)
			if b1 {
				err := SetField(a, v1, p.p2)
				if err != nil {
					panic(err.Error())
				}
			} else {
				panic(errors.New("the type of field name must be string"))
			}
		}

		return WithValues(backup)
	}
}

func SetField(a interface{}, fieldName string, value interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(a))
	if v.CanAddr() == false {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	fieldVal := v.FieldByName(fieldName)
	if fieldVal.CanSet() {
		fieldVal.Set(reflect.ValueOf(value))
	}

	return nil
}
