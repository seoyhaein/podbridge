package podbridge

import (
	"testing"
)

func TestRemoveContainerId(t *testing.T) {
	LC = InitLc()

	LC.AddImagesId("7445a9646150")
	LC.AddImagesId("1b5c7b6fdac0")

	// pod
	LC.AddPodId("test")
	LC.AddContainerInPod("test", "123")
	LC.AddContainerInPod("test", "1111111111111111111111")
	LC.AddContainerInPod("bbbb", "12121")

	LC.ToYaml()
}
