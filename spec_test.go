package podbridge

import (
	"fmt"
	"testing"

	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"
)

// TODO
// Test SetField
// 나머지 필드들에 대한 테스트를 진행해야한다.
// 하지만 세부적인 테스트는 개발진행하면서 진행한다.

func TestSetField(t *testing.T) {

	spec := new(specgen.SpecGenerator)
	fmt.Println("01: 테스트전")
	fmt.Println(spec.Name)
	SetField(spec, "Name", "test")
	fmt.Println("01: 테스트후 test 가 나오면 정상")
	fmt.Println(spec.Name)

	fmt.Println("02: 테스트전")

	for _, ss := range spec.Command {
		fmt.Println(ss)
	}
	ss := []string{"hello", "world"}

	SetField(spec, "Command", ss)
	fmt.Println("02: 테스트후 hello world 가 나오면 정상")

	for _, ss := range spec.Command {
		fmt.Println(ss)
	}
}

// TODO
// Test WithValues

func TestWithValues(t *testing.T) {
	ps := new(pair)
	ps1 := new(pair)

	ps.p1 = "Name"
	ps.p2 = "hello world"

	ps1.p1 = "Image"
	ps1.p2 = "ubuntu"

	opt := WithValues(ps, ps1)
	old := opt(Spec)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("hello world 가 나오면 정상")

	fmt.Println("Spec.Image", Spec.Image)
	fmt.Println("ubuntu 가 나오면 정상")

	old(nil)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("아무것도 안나오면 정상")

	fmt.Println("Spec.Image", Spec.Image)
	fmt.Println("alpine:latest 가 나오면 정상")
}

// nil 값은 복사가 안되고 원본 그대로 나온다.
func TestDeepCopy(t *testing.T) {

	Spec1 := new(specgen.SpecGenerator)
	Spec2 := new(specgen.SpecGenerator)

	Spec1.Name = "babo"

	//deepcopy.DeepCopy(Spec2, Spec1)

	//fmt.Println(Spec2.Name)

	deepcopy.DeepCopy(Spec1, Spec2)
	fmt.Println(Spec1.Name)
}
