package podbridge

import (
	"time"

	"github.com/containers/podman/v4/pkg/specgen"
)

// 컨테이너에 따라 생성된다.
// 따라서 여러개가 생성된다.

type ContainerConfig struct {
	SetSpec                 *bool
	AutoCreateContainerName *bool
}

func (conf *ContainerConfig) TrueSetSpec() *bool {

	conf.SetSpec = PTrue

	return conf.SetSpec
}

func (conf *ContainerConfig) FalseSetSpec() {
	conf.SetSpec = PFalse
	conf.AutoCreateContainerName = PFalse
}

func (conf *ContainerConfig) IsSetSpec() *bool {

	return conf.SetSpec
}

func (conf *ContainerConfig) IsAutoCreateContainerName() *bool {

	return conf.AutoCreateContainerName
}

// 이름을 자동 설정하고 이 메서드를 호출한다.
// 에러 조심하자. nil 의 의미.

func (conf *ContainerConfig) TrueAutoCreateContainerName(spec *specgen.SpecGenerator) *bool {

	// string 이 empty 이면, 즉 세팅이 안되어 있으면
	if IsEmptyString(spec.Name) {
		conf.createSpecContainerName()
		conf.AutoCreateContainerName = PTrue
		return conf.AutoCreateContainerName
	} else { // 만약 Spec.Name 이 세팅되어 있으면 nil 반환.
		if conf.AutoCreateContainerName == PTrue {
			conf.AutoCreateContainerName = PFalse
		}

		return nil
	}
}

func (conf *ContainerConfig) FalseAutoCreateContainerName() {
	conf.AutoCreateContainerName = PFalse
}

// TODO apis.go 로 이동 및 옵션을 만들어서 이름을 자동으로 만들어 줄지 설정할 수 있도록 한다.
// 일단 최초 컨테이너가 생성된 시점의 시간을 기록한다.
// 추가적으로 기록될 필요가 있는 정보가 있으면 추가한다.

func (conf *ContainerConfig) createSpecContainerName() {
	Spec.Name = time.Now().Format("20220702-15h04m05s")
}
