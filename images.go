package podbridge

import (
	"context"
	"errors"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

// TODO  다른 함수로 처리, 오류 있음.
// 함수 및 메서드 정리 필요.

func NewBuildStore() (storage.Store, error) {

	if buildah.InitReexec() {
		return nil, errors.New("buildah init error")
	}

	unshare.MaybeReexecUsingUserNamespace(false)

	buildStoreOptions, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	if err != nil {
		return nil, err
	}

	buildStore, err := storage.GetStore(buildStoreOptions)
	if err != nil {
		return nil, err
	}

	return buildStore, nil

}

func SetFromImage(fromImage string) *buildah.BuilderOptions {

	if IsEmptyString(fromImage) {
		return nil
	}

	return &buildah.BuilderOptions{
		FromImage: fromImage,
	}
}

func NewBuilder(ctx context.Context, store storage.Store, options *buildah.BuilderOptions) (*buildah.Builder, error) {

	builder, err := buildah.NewBuilder(ctx, store, *options)

	return builder, err
}

// 특수한 용도로만 사용된다.
// repository 는 이미지 이름 ex> docker.io/busybox
// TODO 시나리오를 생각하자.
// 임시로 일단 이렇게 하자.

func BuildCustomImage(ctx context.Context, builder *buildah.Builder, store storage.Store, repository string) (*string, error) {

	imageRef, err := is.Transport.ParseStoreReference(store, repository)
	if err != nil {
		return nil, err
	}

	imageId, _, _, err := builder.Commit(ctx, imageRef, buildah.CommitOptions{})
	if err != nil {
		return nil, err
	}

	return &imageId, nil

}

// 마지막에 항상 호출 해줘야 함.

func DeleteAndShutdown(store storage.Store, builder *buildah.Builder) error {

	_, err := store.Shutdown(false)

	if err != nil {
		return err
	}

	err = builder.Delete()

	if err != nil {
		return err
	}

	return nil
}
