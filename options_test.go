package podbridge

import (
	"fmt"
	"testing"

	"github.com/containers/podman/v4/pkg/specgen"
)

func TestWithBasic(t *testing.T) {

	// name 만 설정 되어 있음.
	basicConfig := InitBasicConfig()

	// WithBasic 에서 내부적으로 spec 에 대한 설정을 한다.
	// 그리고 이것을 넘기기 위해서 신규 spec 을 파라미터로 받는 익명함수를 리턴한다.
	// 이 익명함수에 신규 spec 을 넣어서 설정한 값을 적용시킨다.

	var opt Option
	opt = WithBasic(basicConfig)
	// 이렇게 하면 spec 이 Old 값으로 저장됨으로 에러가 난다.
	//var spec = specgen.NewSpecGenerator(imgName, false)
	var spec = new(specgen.SpecGenerator)
	// TODO 성능 비교 필요. new 로 재활용 하지 않는 것과 재활용하는 경우. 5/28
	// 향수 위의 var spec = new(specgen.SpecGenerator) 은 주석 처리 해야하고 전역적으로 선언된  "Spec   *specgen.SpecGenerator"
	// Spec 을 사용해야함. 이때 Spec = eraseSpec(Spec), 해줘서 지워주고 재활용 해야함. TODO default 적용을 해주는 문제도 생각해보자. 다 지우지 말고.
	// 지금 swap 부분도 생각해줘야 한다.

	if opt != nil {
		// 여기서 신규 spec 에 basicConfig 에 적용된  값을 적용 시킨다.
		opt(spec)
		fmt.Println(spec.Name)
		fmt.Println(spec.Image)
	} else {
		fmt.Println("opt nil 이네.")
	}
}

func TestWithTester(t *testing.T) {
	opt := WithTester("hello world")
	var spec = new(specgen.SpecGenerator)
	opt1 := opt(spec)

	opt1(spec)
}
