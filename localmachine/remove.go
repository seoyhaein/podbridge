package localmachine

import (
	"fmt"
	"os/exec"

	pbr "github.com/seoyhaein/podbridge"
)

// TODO 동작확인만 했음. 대충 만듬. 추후 업데이트
// config.go 에서 podbridge.yaml 위치 설정해주는 옵션 설정해주기.

func RemoveContainers() error {
	AllStopContainers()
	// 일단 무식하게 이렇게 한다. 되는지 확인하기.
	for _, id2 := range pbr.Basket.ContainerIds {
		cmd := exec.Command("podman", "rm", "-f", id2)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Errorf("cannot find process")

		}
		fmt.Printf("Status is: %s", string(out))

	}
	return nil
}

func AllStopContainers() error {
	// 모든 컨테이너 중지
	for _, id := range pbr.Basket.ContainerIds {
		cmd := exec.Command("podman", "stop", id)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Errorf("cannot find process")

		}
		fmt.Printf("Status is: %s", string(out))
	}

	return nil
}

func RemovePods() error {
	return nil
}

// TODO buildah 로 빌드된것 안지워지는 경우가 있음.

func RemoveImages() error {
	RemoveContainers()

	for _, id := range pbr.Basket.ImageIds {
		cmd := exec.Command("podman", "rmi", id)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Errorf("cannot find process")

		}
		fmt.Printf("Status is: %s", string(out))
	}

	return nil
}

func RemoveVolumes(lc *pbr.ListCreated) error {
	return nil
}
