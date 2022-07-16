package podbridge

import (
	"context"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

func init() {
	if buildah.InitReexec() {
		return
	}
	unshare.MaybeReexecUsingUserNamespace(false)
}

func newBuildStore() (storage.Store, error) {
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

func setFromImage(fromImage string) *buildah.BuilderOptions {

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

func buildCustomImage(ctx context.Context, builder *buildah.Builder, store storage.Store, repository string) (*string, error) {

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
