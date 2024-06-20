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
	"log/slog"
	"time"

	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/pkg/cleaner"
	"github.com/B1NARY-GR0UP/openalysis/pkg/marker"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/robfig/cron/v3"
)

// TODO: support hooks (pre-handle, post-handle)

var ErrReachedRetryTimes = errors.New("error reached retry times")

var (
	GlobalCleaner = cleaner.New()
	GlobalMarker  = marker.New()
)

func Start(ctx context.Context) error {
	slog.Info("openalysis service started")

	errC := make(chan error, 1)

	if err := GlobalCleaner.AddStrategies(config.GlobalConfig.Cleaner...); err != nil {
		return err
	}
	if err := GlobalMarker.AddStrategies(config.GlobalConfig.Marker...); err != nil {
		return err
	}

	c := cron.New()
	if err := AddCronFunc(ctx, c, errC); err != nil {
		return err
	}

	slog.Info("init task starts now")
	startInit := time.Now()
	tx := storage.DB.Begin()
	// if init failed, stop service
	if err := InitTask(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	slog.Info("init task completed", "time", time.Since(startInit).String())

	c.Start()
	defer c.Stop()

	if err := util.WaitSignal(errC); err != nil {
		slog.Error("receive error signal", "err", err.Error())
	}

	slog.Info("openalysis service stopped")
	return nil
}

func Restart(ctx context.Context) error {
	slog.Info("openalysis service restarted")

	errC := make(chan error, 1)

	if err := GlobalCleaner.AddStrategies(config.GlobalConfig.Cleaner...); err != nil {
		return err
	}
	if err := GlobalMarker.AddStrategies(config.GlobalConfig.Marker...); err != nil {
		return err
	}

	c := cron.New()
	if err := AddCronFunc(ctx, c, errC); err != nil {
		return err
	}

	c.Start()
	defer c.Stop()

	if err := util.WaitSignal(errC); err != nil {
		slog.Error("receive error signal", "err", err.Error())
	}

	slog.Info("openalysis service stopped")
	return nil
}

func AddCronFunc(ctx context.Context, c *cron.Cron, errC chan error) error {
	if _, err := c.AddFunc(config.GlobalConfig.Backend.Cron, func() {
		slog.Info("update task starts now")
		startUpdate := time.Now()
		i := 0
		for {
			i++
			tx := storage.DB.Begin()
			err := UpdateTask(ctx, tx)
			if err == nil {
				tx.Commit()
				break
			}
			slog.Error("error doing update task", "err", err.Error())
			tx.Rollback()
			if i == config.GlobalConfig.Backend.Retry {
				errC <- ErrReachedRetryTimes
				break
			}
			slog.Info("transaction rollback and retry")
		}
		slog.Info("update task completed", "time", time.Since(startUpdate).String())
	}); err != nil {
		return err
	}
	return nil
}
