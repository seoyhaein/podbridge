package podbridge

import (
	"os"

	"github.com/seoyhaein/utils"
	"gopkg.in/yaml.v3"
)

type ListCreated struct {
	ImageIds     []string `yaml:"Images,flow"`
	ContainerIds []string `yaml:"Containers,flow"`
	PodIds       []string `yaml:"Pods,flow"`
	VolumeNames  []string `yaml:"Volumes,flow"`
}

var (
	LC            *ListCreated
	podbridgePath = "podbridge.yaml"
)

func init() {
	// TODO 추후 살펴보자
	LC = InitLc()
}

//ToYaml output to yaml file
func (lc *ListCreated) ToYaml() {
	d, err := yaml.Marshal(lc)

	if err != nil {
		return
	}
	// 기존이 있는 파일을 덥어 쒸운다.
	f, err := os.Create("podbridge.yaml")
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return
	}

	f.Write(d)
	f.Sync()
}

// Deprecated: Not used, but left for now.
func (lc *ListCreated) ToListCreated() *ListCreated {

	temp, err := toListCreated()
	if err != nil {
		return nil
	}

	r := appendListCreated(lc, temp)
	if r == nil {
		return nil
	}

	return r
}

func (lc *ListCreated) AddImagesId(imgId string) *ListCreated {

	//TODO temp 적용이 잘됬는지 확인한다.
	var temp *ListCreated
	lc.ImageIds = append(lc.ImageIds, imgId)

	temp = lc
	return temp
}

func (lc *ListCreated) AddContainerId(containerId string) *ListCreated {

	//TODO temp 적용이 잘됬는지 확인한다.
	var temp *ListCreated
	lc.ContainerIds = append(lc.ContainerIds, containerId)

	temp = lc
	return temp
}

func (lc *ListCreated) AddPodId(podId string) *ListCreated {

	//TODO temp 적용이 잘됬는지 확인한다.
	var temp *ListCreated
	lc.PodIds = append(lc.PodIds, podId)

	temp = lc
	return temp
}

func (lc *ListCreated) AddVolumeName(volumeName string) *ListCreated {

	//TODO temp 적용이 잘됬는지 확인한다.
	var temp *ListCreated
	lc.VolumeNames = append(lc.VolumeNames, volumeName)

	temp = lc
	return temp
}

//toListCreated convert the contents of the podbridge.yaml file to ListCreated
func toListCreated() (*ListCreated, error) {
	var (
		err   error
		bytes []byte
		b     bool
		lc    *ListCreated
	)
	lc = new(ListCreated)
	// 파일이 없을때
	if b, err = utils.FileExists(podbridgePath); b == false {
		f := createPodbridgeYaml()
		if f == nil {
			return nil, err
		}

		return lc, nil
	}

	if bytes, err = os.ReadFile(podbridgePath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, lc); err != nil {
		return nil, err
	}

	return lc, nil

}

// Deprecated: Not used, but left for now.
func appendListCreated(src *ListCreated, temp *ListCreated) *ListCreated {

	if src == nil || temp == nil {
		return nil
	}

	tImages := len(temp.ImageIds)
	tContainers := len(temp.ContainerIds)
	tPods := len(temp.PodIds)
	tVols := len(temp.VolumeNames)

	if tImages > 0 {
		for _, i := range temp.ImageIds {
			src.ImageIds = append(src.ImageIds, i)
		}
	}

	if tContainers > 0 {
		for _, c := range temp.ContainerIds {
			src.ContainerIds = append(src.ContainerIds, c)
		}
	}

	if tPods > 0 {
		for _, p := range temp.PodIds {
			src.PodIds = append(src.PodIds, p)
		}
	}

	if tVols > 0 {
		for _, v := range temp.VolumeNames {
			src.VolumeNames = append(src.VolumeNames, v)
		}
	}

	return src
}

//createPodbridgeYaml create podbridge.yaml
func createPodbridgeYaml() *os.File {
	var (
		f   *os.File
		err error
	)

	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	f, err = os.Create("podbridge.yaml")
	if err != nil {
		return nil
	}
	return f
}

//InitLc used only in the init() function.
func InitLc() *ListCreated {
	temp, err := toListCreated()

	if err != nil {
		return nil
	}

	return temp
}
