package podbridge

import (
	"fmt"
	"testing"

	"github.com/containers/podman/v4/pkg/specgen"
)

// TODO 차후 신규 테스트 진행 필요
func TestWithBasic(t *testing.T) {

	// name 만 설정 되어 있음.
	basicConfig := InitBasicConfig()

	// opt(여기에는 항상 Spec 이 들어가야 한다.)
	var opt Option
	opt = WithBasicConfig(basicConfig)
	// 이렇게 하면 spec 이 Old 값으로 저장됨으로 에러가 난다.
	//var spec = specgen.NewSpecGenerator(imgName, false)
	//var spec = new(specgen.SpecGenerator)
	// TODO 성능 비교 필요. new 로 재활용 하지 않는 것과 재활용하는 경우. 5/28
	// 향수 위의 var spec = new(specgen.SpecGenerator) 은 주석 처리 해야하고 전역적으로 선언된  "Spec   *specgen.SpecGenerator"
	// Spec 을 사용해야함. 이때 Spec = eraseSpec(Spec), 해줘서 지워주고 재활용 해야함. TODO default 적용을 해주는 문제도 생각해보자. 다 지우지 말고.
	// 지금 swap 부분도 생각해줘야 한다.

	if opt != nil {
		// 여기서 신규 spec 에 basicConfig 에 적용된  값을 적용 시킨다.
		opt1 := opt(Spec)
		fmt.Println(Spec.Name)
		fmt.Println(Spec.Image)

		opt1(nil)

		fmt.Println("old 나와야 함.")
		fmt.Println(Spec.Name)
		fmt.Println(Spec.Image)

	} else {
		fmt.Println("opt nil 이네.")
	}
}

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

	ps.p1 = "Name"
	ps.p2 = "hello world"

	old := WithValues(Spec, ps)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("hello world 가 나오면 정상")

	old(nil)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("아무것도 안나오면 정상")
}
