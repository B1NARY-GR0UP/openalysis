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
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

func TestInitTask(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()
	err := InitTask(context.Background(), storage.DB) // around 9 min for cloudwego init
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateTask(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()

	for _, group := range config.GlobalConfig.Groups {
		for _, login := range group.Orgs {
			org, err := graphql.QueryOrgInfo(context.Background(), login)
			if err != nil {
				t.Fatal(err)
			}
			repos, err := graphql.QueryRepoNameByOrg(context.Background(), login)
			if err != nil {
				t.Fatal(err)
			}
			cache[org.ID] = repos
		}
	}
	err := UpdateTask(context.Background(), storage.DB)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestart(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()

	err := CachePreheat(context.Background(), storage.DB)
	if err != nil {
		t.Fatal(err)
	}

	slog.Info("update task starts now")
	startUpdate := time.Now()
	tx := storage.DB.Begin()
	err = UpdateTask(context.Background(), tx)
	if err == nil {
		tx.Commit()
	} else {
		slog.Error("error doing update task", "err", err.Error())
		tx.Rollback()
		slog.Info("transaction rollback and retry")
	}
	slog.Info("update task completed", "time", time.Since(startUpdate).String())
}

func TestProgressBar(t *testing.T) {
	barOut := progressbar.Default(10, "OUT FOR")
	for range 10 {
		_ = barOut.Add(1)
		barIn := progressbar.Default(10, "IN FOR")
		for range 10 {
			_ = barIn.Add(1)
			time.Sleep(time.Second * 1)
		}
	}
}

func TestTransaction(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()
	operation := func(ctx context.Context, db *gorm.DB, count int) error {
		_ = storage.CreateGroup(ctx, db, &model.Group{
			Name:             "test",
			IssueCount:       999,
			PullRequestCount: 999,
			StarCount:        999,
			ForkCount:        999,
			ContributorCount: 999,
		})
		if count == 15 {
			return nil
		}
		return errors.New("error test transaction")
	}
	i := 0
	for {
		tx := storage.DB.Begin()
		err := operation(context.Background(), tx, i)
		if err == nil {
			tx.Commit()
			log.Println("transaction commit")
			break
		}
		tx.Rollback()
		log.Println("transaction rollback")
		i++
		time.Sleep(time.Second * 1)
	}
}

func TestCachePreheat(t *testing.T) {
	config.GlobalConfig.ReadInConfig("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()

	latest := make(map[string][]string)
	for _, group := range config.GlobalConfig.Groups {
		for _, login := range group.Orgs {
			org, err := graphql.QueryOrgInfo(context.Background(), login)
			if err != nil {
				t.Fatal(err)
			}
			repos, err := graphql.QueryRepoNameByOrg(context.Background(), login)
			if err != nil {
				t.Fatal(err)
			}
			latest[org.ID] = repos
		}
	}
	err := CachePreheat(context.Background(), storage.DB)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println("latest:", latest)
	for k, v := range latest {
		fmt.Println(k)
		fmt.Println(len(v))
	}
	//fmt.Println("cache:", cache)
	fmt.Println()
	for k, v := range cache {
		fmt.Println(k)
		fmt.Println(len(v))
	}
	fmt.Println(compareMaps(latest, cache))
}

func compareMaps(map1, map2 map[string][]string) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		if value2, ok := map2[key]; !ok {
			return false
		} else {
			sort.Strings(value1)
			sort.Strings(value2)
			if !reflect.DeepEqual(value1, value2) {
				return false
			}
		}
	}

	return true
}

func TestMap(t *testing.T) {
	map1 := map[string][]string{
		"hello": []string{"1", "3", "5"},
		"world": []string{"2", "4", "6"},
	}
	map2 := map[string][]string{
		"hello": []string{"3", "1", "5"},
		"world": []string{"2", "4", "6"},
	}
	fmt.Println(compareMaps(map1, map2))
}
