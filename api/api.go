package api

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/cron"
	"github.com/B1NARY-GR0UP/openalysis/storage"
)

// TODO: AddGroups SetDataSource SetBackend SetCron SetToken

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
