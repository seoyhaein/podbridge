package podbridge

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/containers/podman/v4/pkg/bindings/containers"
)

// TODO 테스트 진행해야함. 오류 있음.

func TestStartContainerWithSpec(t *testing.T) {

	var (
		finally  Option
		finally1 Option
	)

	sockDir := DefaultLinuxSockDir()
	ctx, err := NewConnection(sockDir, context.Background())

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

	/*tspec := specgen.NewSpecGenerator(busybox, false)
	tspec.Terminal = false
	tspec.Name = "hellobabo"*/

	if conf.AutoCreateContainerName == PFalse || conf.AutoCreateContainerName == nil { // 설정되어 있으면
		name := new(pair)
		name.p1 = "Name"
		name.p2 = time.Now().Format(time.RFC3339)

		opt1 := WithValues(name)
		finally1 = opt1(Spec)
		opt1(Spec)
	}

	//b := reflect.DeepEqual(tspec, Spec)

	/*if b == false {
		fmt.Println("not same")
	} else {
		fmt.Println("same")
	}*/

	b := conf.TrueSetSpec()

	if b == PTrue {
		_, err := containers.CreateWithSpec(*ctx, Spec, &containers.CreateOptions{})
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("Creating %s container using %s image...\n", Spec.Name, Spec.Image)

		//result := StartContainerWithSpec(ctx, conf)

		Finally(finally)

		if conf.AutoCreateContainerName == PFalse || conf.AutoCreateContainerName == nil {
			Finally(finally1)
		}

	}

}
