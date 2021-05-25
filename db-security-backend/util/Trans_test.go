package util

import (
	"fmt"
	"testing"
)

func TestTrans(t *testing.T) {
	var input = []int{1, 1, 1, 1, 1, 1, 0, 1, 0}
	voting := MajorityVoting(input, 3)
	fmt.Println(voting)
}
