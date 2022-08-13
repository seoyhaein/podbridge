package podbridge

import (
	"testing"
)

func TestRemoveContainerId(t *testing.T) {
	basket, _ := MustFirstCall()

	basket.AddImagesId("7445a9646150")
	basket.AddImagesId("1b5c7b6fdac0")

	// pod
	basket.AddPodId("test")
	basket.AddContainerInPod("test", "123")
	basket.AddContainerInPod("test", "1111111111111111111111")
	basket.AddContainerInPod("bbbb", "12121")

	// 여러번 호출해도 문제없는지 테스트 해야한다. 즉, podbridge.yaml 이 지속적으로 업데이트 되는지.
	basket.Save()
	Reset()
}
