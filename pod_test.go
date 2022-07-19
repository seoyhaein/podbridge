package podbridge

import (
	"context"
	"fmt"
	"testing"
)

func TestPodWithSpec(t *testing.T) {

	sockDir := defaultLinuxSockDir()
	ctx, err := NewConnection(context.Background(), sockDir)

	if err != nil {
		fmt.Println("error")
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
}
