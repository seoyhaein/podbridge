package podbridge

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/seoyhaein/utils"
)

func TestStartContainerWithSpec(t *testing.T) {

	var (
		finally  Option
		finally1 Option
	)

	ctx, err := NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error")
	}

	conf := new(ContainerConfig)
	image := new(pair)
	terminal := new(pair)

	busybox := "docker.io/busybox"

	terminal.p1 = "Terminal"
	terminal.p2 = false

	image.p1 = "Image"
	image.p2 = busybox

	opt := WithValues(image, terminal)
	finally = opt(Spec)
	opt(Spec)

	conf.TrueAutoCreateContainerName(Spec)

	if conf.AutoCreateContainerName == utils.PFalse || conf.AutoCreateContainerName == nil { // 설정되어 있으면
		name := new(pair)
		name.p1 = "Name"
		name.p2 = time.Now().Format("20220702-15h04m05s")

		opt1 := WithValues(name)
		finally1 = opt1(Spec)
		opt1(Spec)
	}

	b := conf.TrueSetSpec()

	if b == utils.PTrue {

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result := ContainerWithSpec(ctx, conf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s \n", result.ID, result.Name)

			for i, s := range result.Warnings {
				fmt.Printf("warning(%d): %s \n", i, s)
			}
		}

		Finally(finally)

		if conf.AutoCreateContainerName == utils.PFalse || conf.AutoCreateContainerName == nil {
			Finally(finally1)
		}

	}

}

func TestSetFieldVolume(t *testing.T) {

	var (
		finally  Option
		finally1 Option
	)

	ctx, err := NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error")
	}

	conf := new(ContainerConfig)
	image := new(pair)
	terminal := new(pair)

	busybox := "docker.io/busybox"

	terminal.p1 = "Terminal"
	terminal.p2 = false

	image.p1 = "Image"
	image.p2 = busybox

	vol := new(VolumeConfig)
	volconfig := vol.genVolumeCreateOptions("govol", "local", "/opt")
	named, erro := CreateNamedVolume(ctx, volconfig)

	if erro == nil {
		fmt.Println("에러")
	}

	volumes := new(pair)
	volumes.p1 = "Volumes"
	volumes.p2 = named

	opt := WithValues(image, terminal, volumes)
	finally = opt(Spec)

	conf.TrueAutoCreateContainerName(Spec)

	//	if conf.AutoCreateContainerName == PFalse || conf.AutoCreateContainerName == nil { // 설정되어 있으면
	//		name := new(pair)
	//		name.p1 = "Name"
	//		name.p2 = time.Now().Format("20220702-15h04m05s")
	//
	//		opt1 := WithValues(name)
	//		finally1 = opt1(Spec)
	//		opt1(Spec)
	//	}

	b := conf.TrueSetSpec()

	if b == utils.PTrue {

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result := ContainerWithSpec(ctx, conf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s \n", result.ID, result.Name)

			for i, s := range result.Warnings {
				fmt.Printf("warning(%d): %s \n", i, s)
			}
		}

		Finally(finally)

		if conf.AutoCreateContainerName == utils.PFalse || conf.AutoCreateContainerName == nil {
			Finally(finally1)
		}

	}

}

// TODO 오류 있음.
/*func TestSetVolume(t *testing.T) {

	sockDir := defaultLinuxSockDir()
	ctx, err := NewConnection(context.Background(), sockDir)

	if err != nil {
		fmt.Println("error")
	}

	vol := new(VolumeConfig)
	volconfig := vol.genVolumeCreateOptions("govol", "local", "/opt")

	named, erro := CreateNamedVolume(ctx, volconfig)

	if erro != nil {
		fmt.Println("에러")
	}

	volumes := new(pair)
	volumes.p1 = "Volumes"
	volumes.p2 = named

	opt := WithValues(volumes)
	finally := opt(Spec)

	Finally(finally)
}*/

// TODO volume 연결해서 컨테이너 만들기 테스트

// pod
// TODO 공통적인 코드는 따로 함수로 만들어서 테스트에서 사용하자.

func TestPodSet(t *testing.T) {
	var (
		finally  Option
		finally1 Option
	)

	ctx, err := NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error")
	}

	conf := new(ContainerConfig)
	podConf := new(PodConfig)

	image := new(pair)
	terminal := new(pair)
	pod := new(pair)

	busybox := "docker.io/busybox"

	terminal.p1 = "Terminal"
	terminal.p2 = false

	image.p1 = "Image"
	image.p2 = busybox

	conf.TrueAutoCreateContainerName(Spec)
	podConf.TrueAutoCreatePodNameAndHost(PodSpec)
	b := podConf.TrueSetPodSpec()

	if b == nil || b == utils.PFalse {
		fmt.Println("failed")
		return
	}
	result := PodWithSpec(ctx, podConf)

	if result.success == false || result == nil {
		fmt.Println("failed")
		return
	}

	pod.p1 = "Pod"
	pod.p2 = result.GetPodId()

	opt := WithValues(image, terminal, pod)
	finally = opt(Spec)
	opt(Spec)

	if conf.AutoCreateContainerName == utils.PFalse || conf.AutoCreateContainerName == nil { // 설정되어 있으면
		name := new(pair)
		name.p1 = "Name"
		name.p2 = time.Now().Format("20220702-15h04m05s")

		opt1 := WithValues(name)
		finally1 = opt1(Spec)
		opt1(Spec)
	}

	bp := conf.TrueSetSpec()

	if bp == utils.PTrue {

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		result := ContainerWithSpec(ctx, conf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s \n", result.ID, result.Name)

			for i, s := range result.Warnings {
				fmt.Printf("warning(%d): %s \n", i, s)
			}
		}

		Finally(finally)

		if conf.AutoCreateContainerName == utils.PFalse || conf.AutoCreateContainerName == nil {
			Finally(finally1)
		}

	}
}

func TestSampleContainer(t *testing.T) {

}
