package podbridge

import (
	"strings"
)

var (
	pTrue = true
	PTrue = &pTrue

	pFalse = false
	PFalse = &pFalse
)

// string 이 empty 인경우 true, string 일 경우 false 이다.

func IsEmptyString(s string) bool {

	r := len(strings.TrimSpace(s))

	if r == 0 {
		return true
	}
	return false
}
