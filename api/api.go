package api

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/db"
)

// TODO: main 应该只负责 oa 的初始化以及使用，不负责数据库初始化，配置文件读取等
// TODO: 配置文件读取，数据库读取，开始服务器等都应该在 api 层提供

func Start(path string) {
	Init(path)
	// TODO
}

func Init(path string) {
	config.Init(path)
	db.Init()
	// NOTE: graphql client MUST initialize before rest client due to dependency
	graphql.Init()
	rest.Init()
}

func AddGroups(groups ...config.Group) {
	fmt.Println(groups)
}
