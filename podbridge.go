package podbridge

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/containers/podman/v4/pkg/specgen"
)

func CreatePod(ctx *context.Context, podId string) error {
	var podExistsOptions pods.ExistsOptions

	podExists, err := pods.Exists(*ctx, podId, &podExistsOptions)

	if err != nil {
		return err
	}
	// 기존에 pod 가 존재할 경우는 무조건 해당 pod 를 지운다.
	if podExists {
		return nil
	}

	return nil
}

func CreateContainer(ctx *context.Context, conName string, imgName string) (containerID string) {
	// types.go 참고.

	var containerExistsOptions containers.ExistsOptions

	containerExistsOptions.External = PFalse
	containerExists, err := containers.Exists(*ctx, conName, &containerExistsOptions)
	if err != nil {
		//return nil
	}

	if containerExists {
		var containerInspectOptions containers.InspectOptions
		containerInspectOptions.Size = PFalse
		ins, err := containers.Inspect(*ctx, conName, &containerInspectOptions)
		if err != nil {
			//log.Fatalln(err)
		}

		if ins.State.Running {
			fmt.Printf("%s container already running", conName)
		} else {
			containerID = ins.ID
		}
	} else {
		imageExists, err := images.Exists(*ctx, conName, nil)
		if err != nil {
			//log.Fatalln(err)
		}

		if !imageExists {
			_, err := images.Pull(*ctx, conName, nil)
			if err != nil {
				//log.Fatalln(err)
			}
		}

		fmt.Printf("Creating %s container using %s image...\n", conName, imgName)

		spec := CreateContainerSpec(imgName)
		/*s := specgen.NewSpecGenerator(imgName, false)
		s.Name = conName*/

		createResponse, err := containers.CreateWithSpec(*ctx, spec, nil)
		if err != nil {
			//log.Fatalln(err)
		}

		containerID = createResponse.ID
	}

	return
}

func CreateContainerSpec(imgName string, options ...Option) *specgen.SpecGenerator {

	var spec = specgen.NewSpecGenerator(imgName, false)

	for _, opt := range options {
		if opt != nil {
			opt(spec)
		}
	}

	return spec
}

func InitBasicConfig() *BasicConfig {
	return &BasicConfig{
		Name: "hello",
	}
}
