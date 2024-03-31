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

package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
)

func TestQueryContributorCountByOrg(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	Init()
	count, err := QueryContributorCountByOrg(context.Background(), DB, "O_kgDOCEYWXQ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(count)
}

func TestQueryContributorCountByGroup(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	Init()
	count, err := QueryContributorCountByGroup(context.Background(), DB, "cloudwego")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(count)
}

func TestFor(t *testing.T) {
	var sli []model.Group
	for _, group := range sli {
		fmt.Println("group: ", group)
	}
}

func TestCreate(t *testing.T) {
	err := config.GlobalConfig.ReadInConfig("../default.yaml")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = Init()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = DB.Create([]model.Group{
		{
			Name: "init",
		},
	}).Error
	if err != nil {
		t.Fatal(err.Error())
	}
}
