package main

import (
	"github.com/B1NARY-GR0UP/openalysis/cmd"
)

// TODO
// 每次 api query 都需要插入的 model: repo, contributor
// 需要使用 cursor 增量更新的 model: issue, pr
//
// 1. 读取配置文件
// 2. 初始化数据库
// 3. 插入 groups 表
// 4. 根据配置拉取 api 插入 orgs
// 5. 获取 orgs 下的所有 repos
// 6. 插入 repos 表
// 7. 第一次需要全部拉取 issue 和 pr （只有 OPEN 的 issue 和 pr 才需要分析）
// 8. 插入 groups, orgs, repos, issues, prs, contributors 表
//
// 通过定时任务插入 repo 和 contributor 条目，使用 cursor 增量插入 issue 和 pr 条目
// 更新 groups 和 orgs 的 所有 counting 字段
// TODO: repo table 每次更新都需要插入新的 item，但是 NodeID 是一样的，ID 是自增的，结合 CreatedAt 来绘制时间序列图
func main() {
	// TODO: main 应该只负责 oa 的初始化以及使用，不负责数据库初始化，配置文件读取等
	// TODO: 配置文件读取，数据库读取，开始服务器等都应该在 api 层提供
	cmd.Execute()
}
