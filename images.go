package podbridge

import (
	"context"
	"errors"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

// TODO 읽어보자
// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

type PreBuilderOption struct {
	storage.Store
	*buildah.BuilderOptions
	//*buildah.Builder

	ErrorMessage error
}

// TODO 향후 다른 방향으로 생각하자 일단 이것은 남겨두되 사용하지 않는다.

func BeforeMainStartup() {
	if buildah.InitReexec() {
		return
	}
	unshare.MaybeReexecUsingUserNamespace(false)
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
	// TODO 수정해야 함. 다른 옵션들도 담을 수 있는 방향으로 개선해야 함.
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

//TODO Containerfile/Dockerfile 이미지 만드는 함수 제작
