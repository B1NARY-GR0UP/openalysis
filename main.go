package main

import (
	"github.com/B1NARY-GR0UP/openalysis/cmd"
)

func main() {
	// TODO: main 应该只负责 oa 的初始化以及使用，不负责数据库初始化，配置文件读取等
	// TODO: 配置文件读取，数据库读取，开始服务器等都应该在 api 层提供
	cmd.Execute()
}
