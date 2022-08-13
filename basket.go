package podbridge

import (
	"fmt"
	"os"
	"sync"

	"github.com/seoyhaein/utils"
	"gopkg.in/yaml.v3"
)

// Pod 에 넣은 Container 는 Pod 내에 넣는다.
// key containerid, value pod id

type (
	PodInfo struct {
		Id           string   `yaml:"podId"`
		ContainerIds []string `yaml:"containers,flow"`
	}

	ListCreated struct {
		ImageIds     []string   `yaml:"Images,flow"`
		ContainerIds []string   `yaml:"Containers,flow"`
		Pods         []*PodInfo `yaml:"Pods,flow"`
		VolumeNames  []string   `yaml:"Volumes,flow"`

		mutex *sync.Mutex
	}
)

//TODO 중요 LC 는 공유 struct 이므로 race 문제가 발생할 수 있음. 이걸 보완하자.
// mutex ListCreated 에 넣자.
// https://cloudolife.com/2020/04/18/Programming-Language/Golang-Go/Synchronization/Use-sync-Mutex-sync-RWMutex-to-lock-share-data-for-race-condition/
// Basket 은 singleton 이어야함 관련해서 처리해줘야 하고, 지금 은그냥 노출 하는데, api 를 통해서 노출하도록 처리한다.
var (
	Basket        *ListCreated
	podbridgePath = "podbridge.yaml"
	//mutex         = new(sync.Mutex)
)

//MustFirstCall used only in the init() function.
func MustFirstCall() error {
	basket, err := toListCreated()
	Basket = basket
	return err
}

//ToYaml output to yaml file TODO 수정하자.
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
/*func (lc *ListCreated) ToListCreated() *ListCreated {

	temp, err := toListCreated()
	if err != nil {
		return nil
	}

	r := appendListCreated(lc, temp)
	if r == nil {
		return nil
	}

	return r
}*/

// 중복 검사해줘야 한다.

func (lc *ListCreated) AddImagesId(imgId string) *ListCreated {
	r := findImageId(lc, imgId)
	if r == nil || r == utils.PTrue {
		return lc
	}
	lc.mutex.Lock()
	lc.ImageIds = append(lc.ImageIds, imgId)
	lc.mutex.Unlock()
	return lc
}

func (lc *ListCreated) AddContainerId(containerId string) *ListCreated {
	r := findContainerId(lc, containerId)
	if r == nil || r == utils.PTrue {
		return lc
	}
	lc.mutex.Lock()
	lc.ContainerIds = append(lc.ContainerIds, containerId)
	lc.mutex.Unlock()
	return lc
}

func (lc *ListCreated) AddPodId(podId string) *ListCreated {

	r, _ := findPodId(lc, podId)

	if r == nil || r == utils.PTrue {
		return lc
	}

	newPod := &PodInfo{
		Id: podId,
	}
	lc.mutex.Lock()
	lc.Pods = append(lc.Pods, newPod)
	lc.mutex.Unlock()
	return lc
}

func (lc *ListCreated) AddContainerInPod(podId string, containerIds ...string) *ListCreated {
	var newPod *PodInfo

	r, p := findPodId(lc, podId)
	if r == nil {
		return nil
	}
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	// 동일한 podid 가 없으면
	if r == utils.PFalse {
		newPod = &PodInfo{
			Id: podId,
		}
		for _, c := range containerIds {
			newPod.ContainerIds = append(newPod.ContainerIds, c)
		}
		lc.Pods = append(lc.Pods, newPod)
	}
	// 동일한 podid 가 있으면, deepcopy 가 아니므로 상관없다.
	if r == utils.PTrue {
		newPod = p
		n := len(newPod.ContainerIds)
		if n == 0 {
			for _, c := range containerIds {
				newPod.ContainerIds = append(newPod.ContainerIds, c)
			}
			return lc
		}

		// TODO  좀똑똑하게 고치자 향후에..
		var check = true
		for _, c := range containerIds {
			for _, oc := range newPod.ContainerIds {
				if oc == c {
					check = false
				}
			}
		}
		if check {
			for _, c := range containerIds {
				newPod.ContainerIds = append(newPod.ContainerIds, c)
			}
		}
	}
	return lc
}

func (lc *ListCreated) AddVolumeName(volumeName string) *ListCreated {
	r := findVolumeName(lc, volumeName)
	if r == nil || r == utils.PTrue {
		return lc
	}
	lc.mutex.Lock()
	lc.VolumeNames = append(lc.VolumeNames, volumeName)
	lc.mutex.Unlock()
	return lc
}

/*
func (lc *ListCreated) AddPodIdX(podId string, containerIds ...string) *ListCreated {
	if utils.IsEmptyString(podId) {
		return nil
	}

	if lc.PodsX == nil {
		lc.PodsX = make(map[string]string)
	}

	for _, v := range containerIds {
		lc.PodsX[v] = podId
	}

	return lc
}
*/
// 찾으면 true, 못찾으면 false, 에러면 nil
func findImageId(lc *ListCreated, imageId string) *bool {
	if lc == nil || utils.IsEmptyString(imageId) {
		return nil
	}
	for _, id := range lc.ImageIds {
		if id == imageId {
			return utils.PTrue
		}
	}
	return utils.PFalse
}

// 찾으면 true, 못찾으면 false, 에러면 nil
func findContainerId(lc *ListCreated, conId string) *bool {
	if lc == nil || utils.IsEmptyString(conId) {
		return nil
	}
	for _, id := range lc.ContainerIds {
		if id == conId {
			return utils.PTrue
		}
	}
	return utils.PFalse
}

// 찾으면 true, 못찾으면 false, 에러면 nil
func findVolumeName(lc *ListCreated, name string) *bool {
	if lc == nil || utils.IsEmptyString(name) {
		return nil
	}
	for _, n := range lc.VolumeNames {
		if n == name {
			return utils.PTrue
		}
	}
	return utils.PFalse
}

// 찾으면 true, 못찾으면 false, 에러면 nil
func findPodId(lc *ListCreated, podId string) (*bool, *PodInfo) {
	if lc == nil || utils.IsEmptyString(podId) {
		return nil, nil
	}
	for _, p := range lc.Pods {
		if p.Id == podId {
			return utils.PTrue, p
		}
	}
	return utils.PFalse, nil
}

// 테스트 해보자.
func (lc *ListCreated) RemoveContainerId(containerId string) {
	var index = -1
	for i := 0; i < len(lc.ContainerIds); i++ {
		if lc.ContainerIds[i] == containerId {
			index = i
		}
	}
	if index != -1 {
		lc.mutex.Lock()
		lc.ContainerIds = append(lc.ContainerIds[:index], lc.ContainerIds[index+1:]...)
		lc.mutex.Unlock()
	}
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
	lc.mutex = new(sync.Mutex)
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
/*func appendListCreated(src *ListCreated, temp *ListCreated) *ListCreated {

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
}*/

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

// Reset truncate podbridge.yaml
// ListCreated 를 리셋하는 것은 생각해볼 것
func Reset() error {
	b, err := utils.FileExists(podbridgePath)
	if err != nil {
		return err
	}
	if b {
		err = utils.Truncate(podbridgePath)
		if err != nil {
			return err
		}

	} else {
		// 파일이 없으면
		fmt.Errorf("no file")
	}
	return nil
}
