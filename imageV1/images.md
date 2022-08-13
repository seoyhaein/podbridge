### example

```
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/buildah"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage/pkg/unshare"
	v1 "github.com/seoyhaein/podbridge/imageV1"
)

// Error: could not find slirp4netns, the network namespace can't be configured: exec: "slirp4netns": executable file not found in $PATH
// https://github.com/rootless-containers/slirp4netns

// 컨테이너 에러 발생시 에러 확인
// podman logs (컨테이너 아이디)

// https://stackoverflow.com/questions/49225976/use-sudo-inside-dockerfile-alpine

// alpine
// fallocate -l 10G test.txt
// ls -alh filename

func main() {
	if buildah.InitReexec() {
		return
	}
	unshare.MaybeReexecUsingUserNamespace(false)

	// 여기서 테스트 진행하자.
	opt := v1.NewOption().Other().FromImage("alpine:latest")
	ctx, builder, err := v1.NewBuilder(context.Background(), opt)

	if err != nil {
		fmt.Println("NewBuilder")
		os.Exit(1)
	}
	err = builder.Run("apk update")
	builder.Run("apk add --no-cache bash nano")
	if err != nil {
		fmt.Println("Run1")
		os.Exit(1)
	}

	err = builder.WorkDir("/app")
	if err != nil {
		fmt.Println("WorkDir")
		os.Exit(1)
	}

	err = builder.User("root")
	if err != nil {
		fmt.Println("User")
		os.Exit(1)
	}

	// ADD/Copy 동일함.
	err = builder.Add("./dumpfiles/test.txt", "/app")
	if err != nil {
		fmt.Println("Add")
		os.Exit(1)
	}

	/*err = builder.Run("./file.sh")
	if err != nil {
		fmt.Println("Run2")
		os.Exit(1)
	}*/

	/*err = builder.Run("yarn install --production")
	if err != nil {
		fmt.Println("Run2")
		os.Exit(1)
	}*/

	/*err = builder.Cmd("node", "src/index.js")

	err = builder.Expose("3000")
	if err != nil {
		fmt.Println("Expose")
		os.Exit(1)
	}*/

	sysCtx := &types.SystemContext{}
	image, err := builder.CommitImage(ctx, buildah.Dockerv2ImageManifest, sysCtx, "test01")

	if err != nil {
		fmt.Println("CommitImage")
		os.Exit(1)
	}

	fmt.Println(*image)
}
```

### TODO
- types.SystemContext 관련해서 소스 분석해서 자세한 내용 파악.