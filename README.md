# podbridge

[![Go Reference](https://pkg.go.dev/badge/github.com/seoyhaein/podbridge.svg)](https://pkg.go.dev/github.com/seoyhaein/podbridge)
[![Build Status](https://app.travis-ci.com/seoyhaein/podbridge.svg?branch=main)](https://app.travis-ci.com/seoyhaein/podbridge)

리눅스에서만 구동 가능.
podman 설치 되어야 한다.
## todo
0. containerfile 만들어주는 부분 (이건 일단 후순위로 미룬다.)

1. containerfile 에서 이미지 만들기
1.1 각각 container, volume, image 등의 이름 테그 기타 정보들을 자동으로 규칙적으로 만들어주는 방안 생각해야함.
2. pod 만들기
3. spec 을 설정해서 컨테이너 만들기 (완료)
3.1 컨테이너의 healthcheck, healthcheck 에 따라서 반응하는 루틴 필요
4. 호스트에서 컨테이너로 데이터 전송 podman cp 관련 자료 찾기


## 읽어보기
https://medium.com/safetycultureengineering/an-overview-of-memory-management-in-go-9a72ec7c76a8

## TODO
- ipc 로 restful 연결 임으로 지속적으로 연결을 유지할 필요가 없지 않을까? 
- 여기서는 client 입장에서 접근 한다.
- https://github.com/james-barrow/golang-ipc 참고
- 에러에 관해서 좀 살펴보자.
  참고 : http://cloudrain21.com/golang-graceful-error-handling
- podman/libpod 에서 container.go 잘 살펴보기