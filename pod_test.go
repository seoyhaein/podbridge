package podbridge

import (
	"context"
	"fmt"
	"testing"

	lm "github.com/seoyhaein/podbridge/localmachine"
	"github.com/seoyhaein/utils"
)

func TestPodWithSpec(t *testing.T) {

	//sockDir := defaultLinuxSockDir()
	//ctx, err := NewConnection(context.Background(), sockDir)

	cTx, err := lm.NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error1")
	}

	podConf := new(PodConfig)
	podConf.TrueAutoCreatePodNameAndHost(PodSpec)
	b := podConf.TrueSetPodSpec()

	if b == utils.PTrue {

		result := PodWithSpec(cTx, podConf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s, Hostname: %s \n", result.ID, result.Name, result.Hostname)
		}
	}
}

/*func TestPodWithSpec(t *testing.T) {

	//sockDir := defaultLinuxSockDir()
	//ctx, err := NewConnection(context.Background(), sockDir)

	cTx, err := lm.NewConnectionLinux(context.Background())

	if err != nil {
		fmt.Println("error1")
	}

	podConf := new(PodConfig)
	podConf.TrueAutoCreatePodNameAndHost(PodSpec)
	b := podConf.TrueSetPodSpec()

	if b == utils.PTrue {

		result := PodWithSpec(cTx, podConf)

		if result.success {
			fmt.Printf("ID: %s, Name: %s, Hostname: %s \n", result.ID, result.Name, result.Hostname)
		}
	}
	fmt.Println("hello")
}*/
