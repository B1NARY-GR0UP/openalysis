package api

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/cron"
	"github.com/B1NARY-GR0UP/openalysis/storage"
)

func Start(ctx context.Context, path string) {
	Init(path)
	cron.Start(ctx)
}

func Restart(ctx context.Context, path string) {
	Init(path)
	cron.Restart(ctx)
}

func Init(path string) {
	config.Init(path)
	storage.Init()
	// NOTE: graphql client MUST initialize before rest client due to dependency
	graphql.Init()
	rest.Init()
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
