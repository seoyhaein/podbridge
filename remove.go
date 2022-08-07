package podbridge

import (
	"os"

	"github.com/seoyhaein/utils"
	"gopkg.in/yaml.v3"
)

type ListCreated struct {
	//TODO *string 으로 할지는 추후 살펴보자.
	ImageIds     []string `yaml:"Images,flow"`
	ContainerIds []string `yaml:"Containers,flow"`
	PodIds       []string `yaml:"Pods,flow"`
	VolumeNames  []string `yaml:"Volumes,flow"`
}

// 처음에만 생성한다.

var (
	LC            *ListCreated
	podbridgePath = "podbridge.yaml"
)

func init() {
	// testing 할때는 주석처리.
	//LC = InitLc()
}

// yaml 로 문서로 출력하자.
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

func InitLc() *ListCreated {
	temp, err := toListCreated()

	if err != nil {
		return nil
	}

	return temp
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

/*
var (
		f  *os.File
		fe error
		d  []byte

		images     int
		containers int
		pods       int
		vols       int

		timages     int
		tcontainers int
		tpods       int
		tvols       int
	)

	images = len(lc.ImageIds)
	containers = len(lc.ContainerIds)
	pods = len(lc.PodIds)
	vols = len(lc.VolumeNames)

	defer func() {
		if fe = f.Close(); fe != nil {
			panic(fe)
		}
	}()
	// 어떻게 해서든 file 은 생성됨.
	// podbridge-store.yaml 가 없다면...
	if _, err := os.Stat("podbridge.yaml"); os.IsNotExist(err) {
		f, fe = os.Create("podbridge.yaml")
		if fe != nil {
			panic(fe)
		}

	} else {
		// podbridge-store.yaml 가 있다면...
		if temp := toListCreated("podbridge-store.yaml"); temp != nil {
			timages = len(temp.ImageIds)
			tcontainers = len(temp.ContainerIds)
			tpods = len(temp.PodIds)
			tvols = len(temp.VolumeNames)
		}

	}
	// 기존 메모리에 데이터가 있으면
	if images > 0 || containers > 0 || pods > 0 || vols > 0 {
		d, err := yaml.Marshal(lc)

		if err != nil {
			return
		}

		// TODO 에러 반환값으로 수정
		if _, fe = f.Write(d); fe != nil {
			return
		}
	}



*/
