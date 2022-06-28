package main

import (
	"context"
	"fmt"
	"reflect"

	pbr "github.com/seoyhaein/podbridge"
)

// 터미널에서 확인하기 위해서 만듬.

func main() {
	sockDir := pbr.DefaultLinuxSockDir()
	_, err := pbr.NewConnection(sockDir, context.Background())

	if err != nil {
		fmt.Println("error")
	}

	findField()
}

func findField() {
	//t := reflect.TypeOf(*pbr.Spec)

	fmt.Println(reflect.TypeOf(*pbr.Spec).Kind().String())
}
