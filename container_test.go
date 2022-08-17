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

// podman inspect contest06 --format '{{.State.Healthcheck}}'

// podman run -dt --name contest07 --health-cmd='HealthCheckTest.sh' --health-interval=0 docker.io/library/test05
// podman inspect contest07 --format '{{.State.Healthcheck}}'
// 직접하면 되는데..흠..

// https://knowledge.broadcom.com/external/article/237006/how-to-custom-health-check-settings-for.html
// https://developers.redhat.com/blog/2019/04/18/monitoring-container-vitality-and-availability-with-podman#what_are_healthchecks_
// https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/building_running_and_managing_containers/assembly_monitoring-containers_building-running-and-managing-containers
// https://devops.stackexchange.com/questions/11501/healthcheck-cmd-vs-cmd-shell
// 안정화 시켜야 한다.
//
func TestContainer01(t *testing.T) {

	cTx, err := NewConnectionLinux(context.Background())
	if err != nil {
		t.Fail()
	}
	// spec 만들기
	conSpec := NewSpec()
	conSpec.SetImage("docker.io/library/test06")

	f := func(spec SpecGen) SpecGen {
		spec.Name = "contest18"
		spec.Terminal = true
		return spec
	}
	conSpec.SetOther(f)
	// 해당 이미지에 해당 shell script 가 있다.
	conSpec.SetHealthChecker("CMD-SHELL /app/healthcheck.sh", "2s", 3, "30s", "1s")
	/*f1 := func(spec SpecGen) SpecGen {
		healthConfig, _ := SetHealthChecker(spec, "CMD-SHELL /app/healthcheck.sh", "2s", 3, "30s", "1s")
		spec.HealthConfig = healthConfig
		return spec
	}*/

	//conSpec.SetOther(f1)

	// container 만들기
	r := CreateContainer(cTx, conSpec)
	fmt.Println("container Id is :", r.ID)
	err = r.Start(cTx)
	r.HealthCheck(cTx, "1s")

}
