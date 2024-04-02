// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"testing"
)

func TestNameWithOwner(t *testing.T) {
	nameWithOwner := "cloudwego/hertz"
	owner, name := SplitNameWithOwner(nameWithOwner)
	fmt.Println("owner:", owner)
	fmt.Println("name:", name)
	res := MergeNameWithOwner(owner, name)
	fmt.Println("res:", res)
}

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
}
