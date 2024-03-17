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

func TestCompareSlices(t *testing.T) {
	fmt.Println(CompareSlices([]int{1, 2, 3}, []int{1, 2, 3}))
	fmt.Println(CompareSlices([]int{1, 2, 3}, []int{1, 2, 3, 4, 5}))
	fmt.Println(CompareSlices([]int{1, 2, 3}, []int{1}))
	fmt.Println(CompareSlices([]int{1, 2, 3, 5}, []int{1, 2, 3, 4}))

	m := map[string]bool{
		"hello": true,
	}
	fmt.Println(m["world"])

}
