package podbridge

import (
	"fmt"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// TODO 공부할게 많다. 힘들다. 젠장.
// 현재 os 의 resource  를 가지고 와서 이것과 비교해서 제한을 둬야 한다.
// mesos resource 참고

// https://access.redhat.com/documentation/ko-kr/red_hat_enterprise_linux/6/html/resource_management_guide/sec-cpu

// cpu share 의 개념
// https://kimmj.github.io/kubernetes/kubernetes-cpu-request-limit/

// 잠깐 조사를 거친 후 드는 생각은 mesos 의 리소스 설정의 경우는 cpu 는 코어 만 적용이되는 듯하다. 이 부분은 실제로 테스트 하면서 많이 조사를 해야 하는 부분이다.
// nomad 참고하자. 이녀석은 specs-go 를 참고하지 않고 자체적으로 struct 를 구현했다. => 표준과 벗어난 거 아닌가?? 흠.

// nomad driver_test.go 내용 참고하자.
/*
var (
	basicResources = &drivers.Resources{
		NomadResources: &structs.AllocatedTaskResources{
			Memory: structs.AllocatedMemoryResources{
				MemoryMB: 256,
			},
			Cpu: structs.AllocatedCpuResources{
				CpuShares: 250,
			},
		},
		LinuxResources: &drivers.LinuxResources{
			CPUShares:        512,
			MemoryLimitBytes: 256 * 1024 * 1024,
		},
	}
)
*/

// 위의 내용을 살펴보면, mesos 에서는 resource 관련해서는 container 와 non-container 의 리소스 설정이 다르다.
// https://mesos.apache.org/documentation/attributes-resources/ -> non-container
// https://mesos.apache.org/documentation/latest/nested-container-and-task-group/ -> container

// TODO error prone!
// return specs.LinuxResource 로 리턴해야 함.

func LimitResources(cores float64, mems float64) *specs.LinuxCPU {

	if cores <= 0 {
		return nil
	}

	LinuxCpus := new(specs.LinuxCPU)

	period, quota := util.CoresToPeriodAndQuota(cores)
	LinuxCpus.Period = &period
	LinuxCpus.Quota = &quota
	LinuxCpus.Cpus = fmt.Sprintf("%f", cores)
	LinuxCpus.Mems = fmt.Sprintf("%f", mems)

	return LinuxCpus
}
