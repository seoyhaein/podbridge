package podbridge

import (
	"testing"
)

func TestRemoveContainerId(t *testing.T) {
	MustFirstCall()

	Basket.AddImagesId("7445a9646150")
	Basket.AddImagesId("1b5c7b6fdac0")

	// pod
	Basket.AddPodId("test")
	Basket.AddContainerInPod("test", "123")
	Basket.AddContainerInPod("test", "1111111111111111111111")
	Basket.AddContainerInPod("bbbb", "12121")

	// 여러번 호출해도 문제없는지 테스트 해야한다. 즉, podbridge.yaml 이 지속적으로 업데이트 되는지.
	Basket.Save()
	Reset()
}
