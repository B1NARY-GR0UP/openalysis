package cron

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/robfig/cron/v3"
	"github.com/schollz/progressbar/v3"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	// 创建一个新的 cron 实例
	c := cron.New()

	// 添加定时任务
	_, err := c.AddFunc("* * * * *", func() {
		fmt.Println("执行定时任务：", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Println("添加定时任务失败:", err)
		return
	}

	// 启动 cron 调度器
	c.Start()

	defer c.Stop()

	// 等待程序运行一段时间以便查看输出
	// 由于是每分钟执行一次，因此可以注释掉此行以使程序一直运行
	time.Sleep(5 * time.Minute)
}

func TestInitTask(t *testing.T) {
	config.Init("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()
	InitTask(context.Background()) // around 9 min for cloudwego init
}

func TestUpdateTask(t *testing.T) {
	config.Init("../default.yaml")
	storage.Init()
	graphql.Init()
	rest.Init()
	cache = make(map[string][]string)
	for _, group := range config.GlobalConfig.Groups {
		for _, login := range group.Orgs {
			repos, err := graphql.QueryRepoNameByOrg(context.Background(), login)
			if err != nil {
				panic("test panic")
			}

			cache[login] = repos
		}
	}
	UpdateTask(context.Background())
}

func TestProgressBar(t *testing.T) {
	barOut := progressbar.Default(10, "OUT FOR")
	for _ = range 10 {
		barOut.Add(1)
		barIn := progressbar.Default(10, "IN FOR")
		for _ = range 10 {
			barIn.Add(1)
			time.Sleep(time.Second * 1)
		}
	}
}
