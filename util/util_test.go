package util

import (
	"fmt"
	"testing"
)

func TestIsEmptySlice(t *testing.T) {
	sli1 := make([]int, 0)
	var sli2 []int
	fmt.Println(IsEmptySlice(sli1))
	fmt.Println(IsEmptySlice(sli2))
}
