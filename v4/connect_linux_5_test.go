package podbridge

import (
	"fmt"
	"os"
	"testing"
)

func TestNewConnection5(t *testing.T) {
	sockDir := fmt.Sprintf("/run/user/%d", os.Getuid())
	fmt.Println(sockDir)
}
