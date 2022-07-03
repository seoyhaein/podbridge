package podbridge

import (
	"context"

	"github.com/containers/podman/pkg/bindings/volumes"
	"github.com/containers/podman/pkg/domain/entities"
)

// TODO 테스트
// 조심해서 살펴보자.

func create() {
	ctx := context.Background()
	// field 의 정보는 아래 링크에서 확인
	// https://docs.podman.io/en/latest/markdown/podman-volume-create.1.html
	volcreateoption := new(entities.VolumeCreateOptions)
	volcreateoption.Name = "my_vol"

	conf := *volcreateoption

	volumes.Create(ctx, conf, &volumes.CreateOptions{})
}
