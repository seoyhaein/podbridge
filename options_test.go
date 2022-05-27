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

	var opt1 Option
	opt1 = WithBasic(basicConfig)
	imgName := "docker.io/centos:latest"

	//var spec = new(specgen.SpecGenerator)
	var spec = specgen.NewSpecGenerator(imgName, false)
	if opt1 != nil {
		opt1(spec)
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
