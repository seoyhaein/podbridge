package podbridge

import (
	"context"
	"fmt"
	"testing"
)

// 이런식으로 하면 될듯한데. 테스트 해보자.
func TestSetOther(t *testing.T) {
	cspec := NewSpec()

	f := func(spec SpecGen) SpecGen {
		spec.Name = "test"
		return spec
	}

	cspec.SetImage("busybox")
	cspec.SetOther(f)

	fmt.Println(cspec.Spec.Image, cspec.Spec.Name)
}

func TestContainer01(t *testing.T) {

	cTx, err := NewConnectionLinux(context.Background())
	if err != nil {
		t.Fail()
	}
	// spec 만들기
	conSpec := NewSpec()
	conSpec.SetImage("docker.io/busybox")

	f := func(spec SpecGen) SpecGen {
		spec.Name = "test02"
		spec.Terminal = true
		return spec
	}
	conSpec.SetOther(f)

	// container 만들기
	r := CreateContainer(cTx, conSpec)
	fmt.Println("container Id is :", r.ID)
	err = r.Start(cTx)

	r.HealthCheck(cTx)
}
