package localmachine

import (
	"testing"

	pbr "github.com/seoyhaein/podbridge"
)

func TestRemove01(t *testing.T) {
	pbr.LC = pbr.InitLc()

	pbr.LC.AddImagesId("7445a9646150")
	pbr.LC.AddImagesId("1b5c7b6fdac0")

	// pod
	pbr.LC.ToYaml()
	RemoveImages()

	pbr.Reset()
}
