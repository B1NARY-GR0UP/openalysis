package api

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/db"
)

// TODO: main 应该只负责 oa 的初始化以及使用，不负责数据库初始化，配置文件读取等
// TODO: 配置文件读取，数据库读取，开始服务器等都应该在 api 层提供

func Start(path string) {
	InitConfig(path)
	InitDB()

}

func InitClient() {
	graphql.Init()
}

// InitDB should execute after InitConfig due to DSN configuration
func InitDB() {
	db.Init()
}

func InitConfig(path string) {
	config.Init(path)
}

func AddGroups(groups ...config.Group) {
	fmt.Println(groups)
}
