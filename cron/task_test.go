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

package cron

import (
	"fmt"
	"testing"
)

type User struct {
	ID    string
	Name  string
	Count int
}

func Test1(t *testing.T) {
	m := make(map[string]User)
	if _, ok := m["1"]; !ok {
		m["1"] = User{ID: "1", Name: "lorain", Count: 0}
	}
	fmt.Println(m) // 0
	if v, ok := m["1"]; ok {
		v.Count += 10
	}
	fmt.Println(m) // 0
}

func Test2(t *testing.T) {
	m := make(map[string]User)
	if _, ok := m["1"]; !ok {
		m["1"] = User{ID: "1", Name: "lorain", Count: 0}
	}
	fmt.Println(m) // 0
	if v, ok := m["1"]; ok {
		v.Count += 10
		m["1"] = v
	}
	fmt.Println(m) // 10
}

func Test3(t *testing.T) {
	m := make(map[string]*User)
	if _, ok := m["1"]; !ok {
		m["1"] = &User{ID: "1", Name: "lorain", Count: 0}
	}
	fmt.Println(m) // 0
	if v, ok := m["1"]; ok {
		v.Count += 10
	}
	fmt.Println(m["1"].Count) // 10
}
