package podbridge

import (
	"context"
	"fmt"
	"testing"
	"time"
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

// podman inspect contest06 --format '{{.State.Healthcheck}}'

// podman run -dt --name contest07 --health-cmd='HealthCheckTest.sh' --health-interval=0 docker.io/library/test05
// podman inspect contest07 --format '{{.State.Healthcheck}}'
// 직접하면 되는데..흠..

// https://knowledge.broadcom.com/external/article/237006/how-to-custom-health-check-settings-for.html
// 안정화 시켜야 한다.
func TestContainer01(t *testing.T) {

	cTx, err := NewConnectionLinux(context.Background())
	if err != nil {
		t.Fail()
	}
	// spec 만들기
	conSpec := NewSpec()
	conSpec.SetImage("docker.io/library/test05")

	f := func(spec SpecGen) SpecGen {
		spec.Name = "contest15"
		spec.Terminal = true
		return spec
	}
	conSpec.SetOther(f)

	f1 := func(spec SpecGen) SpecGen {
		healthConfig, _ := SetHealthChecker(spec, "CMD-SHELL /app/HealthCheckTest.sh", "2s", 3, "30s", "1s")
		spec.HealthConfig = healthConfig
		return spec
	}

	conSpec.SetOther(f1)

	// container 만들기
	r := CreateContainer(cTx, conSpec)
	fmt.Println("container Id is :", r.ID)
	err = r.Start(cTx)
	time.Sleep(time.Second * 2)
	r.HealthCheck(cTx)
}
