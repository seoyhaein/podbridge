package podbridge

import (
	"context"

	"github.com/containers/podman/v4/pkg/bindings/pods"
)

// TODO 구현중
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
