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

package api

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/cron"
	"github.com/B1NARY-GR0UP/openalysis/storage"
)

func Start(ctx context.Context) {
	cron.Start(ctx)
}

func Restart(ctx context.Context) {
	cron.Restart(ctx)
}

func Init(path string) error {
	if err := config.Init(path); err != nil {
		return err
	}
	if err := storage.Init(); err != nil {
		return err
	}
	// NOTE: graphql client MUST initialize before rest client due to dependency
	graphql.Init()
	rest.Init()
	return nil
}

func AddGroups(groups ...config.Group) {
	config.GlobalConfig.Groups = append(config.GlobalConfig.Groups, groups...)
}

func SetDataSource(ds config.DataSource) {
	config.GlobalConfig.DataSource = ds
}

func SetBackend(be config.Backend) {
	config.GlobalConfig.Backend = be
}

func SetCron(spec string) {
	config.GlobalConfig.Backend.Cron = spec
}

func SetToken(token string) {
	config.GlobalConfig.Backend.Token = token
}

func SetRetry(times int) {
	config.GlobalConfig.Backend.Retry = times
}
