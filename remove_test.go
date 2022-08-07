package podbridge

import (
	"testing"
	//"github.com/stretchr/testify/assert"
)

func TestInitLc01(t *testing.T) {

	//assert := assert.New(t)

	// 파일이 없을때와 파일이 있을때 비교하자.
	//lc := InitLc()
	LC = InitLc()
	LC.AddImagesId("7445a9646150")
	LC.AddImagesId("1b5c7b6fdac0")
	LC.ToYaml()

	// 파일이 있을때
	/*n1 := len(lc.ImageIds)
	n2 := len(lc.PodIds)
	n3 := len(lc.VolumeNames)
	n4 := len(lc.ContainerIds)

	assert.Equal(0, n1, "ImageIds")
	assert.Equal(0, n2, "PodIds")
	assert.Equal(0, n3, "VolumeNames")
	assert.Equal(0, n4, "ContainerIds")*/
}

// 기존 파일을 덥어쒸울때
func TestInitLc02(t *testing.T) {
	//assert := assert.New(t)

	// 파일이 없을때와 파일이 있을때 비교하자.
	//lc := InitLc()
	LC = InitLc()
	LC.AddImagesId("6d97ac174ea6")
	LC.ToYaml()

	// 파일이 있을때
	/*n1 := len(lc.ImageIds)
	n2 := len(lc.PodIds)
	n3 := len(lc.VolumeNames)
	n4 := len(lc.ContainerIds)

	assert.Equal(0, n1, "ImageIds")
	assert.Equal(0, n2, "PodIds")
	assert.Equal(0, n3, "VolumeNames")
	assert.Equal(0, n4, "ContainerIds")*/
}
