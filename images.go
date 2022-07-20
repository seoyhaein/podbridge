package podbridge

import (
	"context"
	"errors"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

type PreBuilderOption struct {
	storage.Store
	*buildah.BuilderOptions
	//*buildah.Builder

	ErrorMessage error
}

func NewBuildImage(fromImage string) *PreBuilderOption {

	preBuilderOption := new(PreBuilderOption)

	buildStoreOptions, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())

	if err != nil {
		preBuilderOption.Store = nil
		preBuilderOption.BuilderOptions = nil
		//newBuildImage.Builder = nil
		preBuilderOption.ErrorMessage = err

		return preBuilderOption
	}

	buildStore, err := storage.GetStore(buildStoreOptions)

	if err != nil {
		preBuilderOption.Store = nil
		preBuilderOption.BuilderOptions = nil
		//newBuildImage.Builder = nil
		preBuilderOption.ErrorMessage = err

		return preBuilderOption
	}

	preBuilderOption.Store = buildStore

	if IsEmptyString(fromImage) {
		preBuilderOption.BuilderOptions = nil
		//newBuildImage.Builder = nil
		preBuilderOption.ErrorMessage = errors.New("there is no image name")

		return preBuilderOption
	}
	builderOption := new(buildah.BuilderOptions)
	builderOption.FromImage = fromImage
	preBuilderOption.BuilderOptions = builderOption

	//builder, err := buildah.NewBuilder(ctx, store, *options)

	return preBuilderOption
}

func (pbo *PreBuilderOption) NewBuilder(ctx context.Context) (*buildah.Builder, error) {

	builder, err := buildah.NewBuilder(ctx, pbo.Store, *pbo.BuilderOptions)

	return builder, err
}

func (pbo *PreBuilderOption) DeleteAndShutdown(builder *buildah.Builder) error {
	_, err := pbo.Store.Shutdown(false)

	if err != nil {
		return err
	}

	if builder == nil {
		errors.New("invalid builder, builder is nil ")
	}
	err = builder.Delete()

	if err != nil {
		return err
	}

	return nil
}

func NewBuildStore() (storage.Store, error) {

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
