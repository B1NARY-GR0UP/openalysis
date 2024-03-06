package task

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"log/slog"
)

// TODO: listen to close signal
var closeC chan struct{}

// groups
// key: group name
// value: slice of repos
var groups map[string][]string

// TODO: 记录每个小 task 的执行耗时

func Init() {
	groups = make(map[string][]string)
	for _, group := range config.GlobalConfig.Groups {
		repos := make([]string, 0)
		repos = append(repos, group.Repos...)
		for _, org := range group.Orgs {
			res, err := graphql.QueryRepoNameByOrg(context.Background(), org)
			if err != nil {
				slog.Error("error query repo name by org", "err", err.Error())
				continue
			}
			repos = append(repos, res...)
		}
		groups[group.Name] = repos
	}

	for name, repos := range groups {
		// TODO: 查询每个 repo 的 info
		// TODO: 查询每个 repo 的 contributor 并计算 count
		// TODO: 查询每个 repo 的 issue 和 pr
		// TODO: 插入数据库, contributor, issue, pr 为最细粒度
	}
}
