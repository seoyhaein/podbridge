package podbridge

import (
	"testing"
)

func TestGenExecutorSh(t *testing.T) {
	path := "."
	fileName := "executor.sh"
	cmd := `echo "hello world"`

	_, _, err := genExecutorSh(path, fileName, cmd)
	if err != nil {
		t.Fail()
	}
}
