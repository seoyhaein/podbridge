package podbridge

import (
	"context"
	"fmt"
	"testing"
)

/*func TestPodWithSpec(t *testing.T) {

	//sockDir := defaultLinuxSockDir()
	//ctx, err := NewConnection(context.Background(), sockDir)

	ctx, err := NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error1")
	}

	podConf := new(PodConfig)
	podConf.TrueAutoCreatePodNameAndHost(PodSpec)
	b := podConf.TrueSetPodSpec()

	if b == PTrue {

		result := PodWithSpec(ctx, podConf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s, Hostname: %s \n", result.ID, result.Name, result.Hostname)
		}
	}
}*/

func TestPodWithSpec(t *testing.T) {

	//sockDir := defaultLinuxSockDir()
	//ctx, err := NewConnection(context.Background(), sockDir)

	NewConnectionLinux(context.Background())

	/*if err != nil {
		fmt.Println("error1")
	}

	podConf := new(PodConfig)
	podConf.TrueAutoCreatePodNameAndHost(PodSpec)
	b := podConf.TrueSetPodSpec()

	if b == PTrue {

		result := PodWithSpec(ctx, podConf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s, Hostname: %s \n", result.ID, result.Name, result.Hostname)
		}
	}*/
	fmt.Println("hello")
}
