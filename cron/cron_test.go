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
	"log"
	"log/slog"
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
	_ = config.GlobalConfig.ReadInConfig("../default.yaml")
	_ = storage.Init()
	graphql.Init()
	rest.Init()
	err := InitTask(context.Background(), storage.DB) // around 9 min for cloudwego init
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateTask(t *testing.T) {
	_ = config.GlobalConfig.ReadInConfig("../default.yaml")
	_ = storage.Init()
	graphql.Init()
	rest.Init()

	slog.Info("update task starts now")
	startUpdate := time.Now()
	tx := storage.DB.Begin()
	err := UpdateTask(context.Background(), tx)
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
	_ = config.GlobalConfig.ReadInConfig("../default.yaml")
	_ = storage.Init()
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

func TestClean(t *testing.T) {
	_ = config.GlobalConfig.ReadInConfig("../default.yaml")
	_ = storage.Init()
	ss := []string{
		"`ByteDance` => `TEST PASS`",
		"`蚂蚁` => `Alibaba`",
		"`Beijing` => `Beijing, China`",
	}
	_ = GlobalCleaner.AddStrategies(ss...)
	err := CleanContributorCompanyAndLocation(context.Background(), storage.DB)
	if err != nil {
		t.Fatal(err)
	}
}
