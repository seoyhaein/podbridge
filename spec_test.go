package podbridge

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	deepcopy "github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/containers/podman/v4/pkg/specgen"
)

// TODO
// Test SetField
// 나머지 필드들에 대한 테스트를 진행해야한다.
// 하지만 세부적인 테스트는 개발진행하면서 진행한다.

func TestSetField(t *testing.T) {

	spec := new(specgen.SpecGenerator)
	fmt.Println("01: 테스트전")
	fmt.Println(spec.Name)
	SetField(spec, "Name", "test")
	fmt.Println("01: 테스트후 test 가 나오면 정상")
	fmt.Println(spec.Name)

	fmt.Println("02: 테스트전")

	for _, ss := range spec.Command {
		fmt.Println(ss)
	}
	ss := []string{"hello", "world"}

	SetField(spec, "Command", ss)
	fmt.Println("02: 테스트후 hello world 가 나오면 정상")

	for _, ss := range spec.Command {
		fmt.Println(ss)
	}
}

// TODO
// Test WithValues

func TestWithValues(t *testing.T) {
	ps := new(pair)
	ps1 := new(pair)

	ps.p1 = "Name"
	ps.p2 = "hello world"

	ps1.p1 = "Image"
	ps1.p2 = "ubuntu"

	opt := WithValues(ps, ps1)
	old := opt(Spec)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("hello world 가 나오면 정상")

	fmt.Println("Spec.Image", Spec.Image)
	fmt.Println("ubuntu 가 나오면 정상")

	old(nil)
	fmt.Println("Spec.Name", Spec.Name)
	fmt.Println("아무것도 안나오면 정상")

	fmt.Println("Spec.Image", Spec.Image)
	fmt.Println("alpine:latest 가 나오면 정상")
}

// nil 값은 복사가 안되고 원본 그대로 나온다.
func TestDeepCopy(t *testing.T) {

	Spec1 := new(specgen.SpecGenerator)
	Spec2 := new(specgen.SpecGenerator)

	Spec1.Name = "babo"

	//deepcopy.DeepCopy(Spec2, Spec1)

	//fmt.Println(Spec2.Name)

	deepcopy.DeepCopy(Spec1, Spec2)
	fmt.Println(Spec1.Name)
}

func TestSpecCompare(t *testing.T) {

	sockDir := DefaultLinuxSockDir()
	ctx, err := NewConnection(sockDir, context.Background())

	//centos := "docker.io/centos"
	busybox := "docker.io/busybox"
	tspec := specgen.NewSpecGenerator(busybox, false)
	//tspec.Terminal = true

	ps := new(pair)
	ps.p1 = "Image"
	ps.p2 = busybox

	terminal := new(pair)
	terminal.p1 = "Terminal"
	terminal.p2 = false

	opt := WithValues(ps, terminal)
	old := opt(Spec)

	b := reflect.DeepEqual(tspec, Spec)

	if b == false {
		fmt.Println("not same")
	} else {
		fmt.Println("same")
	}

	_, err = containers.CreateWithSpec(*ctx, Spec, &containers.CreateOptions{})

	if err != nil {
		fmt.Println(err.Error())
	}

	old(nil)

}

func TestTime(t *testing.T) {

	sockDir := DefaultLinuxSockDir()
	ctx, _ := NewConnection(sockDir, context.Background())

	busybox := "docker.io/busybox"
	ps := new(pair)
	ps.p1 = "Image"
	ps.p2 = busybox

	terminal := new(pair)
	terminal.p1 = "Terminal"
	terminal.p2 = false

	name := new(pair)
	name.p1 = "Name"
	//name.p2 = time.Now().Format(time.RFC3339)
	s1 := time.Now().Format(time.RFC3339)
	s4 := time.Now().Format("20220702-15h04m05s")
	fmt.Println(s4)

	// ":" 문자를 "-" 로 다 치환한다.
	s2 := strings.Replace(s1, ":", "-", -1)
	s3 := strings.Replace(s2, "+", "-", -1)

	name.p2 = s4

	fmt.Println(s3)
	opt := WithValues(ps, terminal, name)
	old := opt(Spec)

	_, err := containers.CreateWithSpec(*ctx, Spec, &containers.CreateOptions{})

	if err != nil {
		fmt.Println(err.Error())
	}

	old(nil)
}
