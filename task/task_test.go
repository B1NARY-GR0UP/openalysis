package task

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/robfig/cron/v3"
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

func TestInitRepoTask(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../default.yaml"); err != nil {
		panic(err.Error())
	}
	graphql.Init()
	rest.Init()
	res, err := InitRepoTask("cloudwego/hertz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
