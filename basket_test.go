package podbridge

import (
	"testing"
)

func TestRemoveContainerId(t *testing.T) {
	Basket = InitLc()

	Basket.AddImagesId("7445a9646150")
	Basket.AddImagesId("1b5c7b6fdac0")

	// pod
	Basket.AddPodId("test")
	Basket.AddContainerInPod("test", "123")
	Basket.AddContainerInPod("test", "1111111111111111111111")
	Basket.AddContainerInPod("bbbb", "12121")

	Basket.ToYaml()
	Reset()
}
