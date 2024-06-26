# podbridge

[![Go Reference](https://pkg.go.dev/badge/github.com/seoyhaein/podbridge.svg)](https://pkg.go.dev/github.com/seoyhaein/podbridge)
[![Build Status](https://app.travis-ci.com/seoyhaein/podbridge.svg?branch=main)](https://app.travis-ci.com/seoyhaein/podbridge)
[![CodeFactor](https://www.codefactor.io/repository/github/seoyhaein/podbridge/badge)](https://www.codefactor.io/repository/github/seoyhaein/podbridge)

```
package main

import (
	"context"
	"fmt"

	pbr "github.com/seoyhaein/podbridge"
)

func main() {
	ctx, err := pbr.NewConnectionLinux(context.Background())
	if err != nil {
		panic(err)
	}
	// spec 만들기
	conSpec := pbr.NewSpec()
	conSpec.SetImage("docker.io/library/test07")

	f := func(spec pbr.SpecGen) pbr.SpecGen {
		spec.Name = "container-tester01"
		spec.Terminal = true
		return spec
	}
	conSpec.SetOther(f)
	// 해당 이미지에 해당 shell script 가 있다.
	conSpec.SetHealthChecker("CMD-SHELL /app/healthcheck/healthcheck.sh", "2s", 1, "30s", "1s")

	// container 만들기
	r := pbr.CreateContainer(ctx, conSpec)
	fmt.Println("container Id is :", r.ID)
	result := r.RunT(ctx, "1s")

	v := int(result)
	fmt.Println(v)
}

```

Don't forget !
podman system service -t 0 &

-설치 스크립트 수정 필요. 최신 버전 Makefile 및 메뉴얼 확인 필요. (처리 후 삭제)

