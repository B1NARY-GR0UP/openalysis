package graphql

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
	"time"
)

func TestQueryRepoInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	info, err := QueryRepoInfo(context.Background(), "cloudwego", "hertz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestQueryOrgInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	info, err := QueryOrgInfo(context.Background(), "cloudwego")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestQueryRepoNameByOrg(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	res, err := QueryRepoNameByOrg(context.Background(), "cloudwego")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestQueryIssueInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	issues, lastUpdate, err := QueryIssueInfoByRepo(context.Background(), "cloudwego", "hertz", time.Time{})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("length:", len(issues))
	fmt.Println(lastUpdate.String())
	for _, issue := range issues {
		fmt.Println(issue)
		fmt.Println(len(issue.Assignees.Nodes)) // all alloc memory
	}

	//lastUpdate := time.Now().UTC()
	//time.Sleep(time.Second * 30)
	//issues, _, err := QueryIssueInfo(context.Background(), "rainiring", "test", lastUpdate)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(len(issues))
	//for _, issue := range issues {
	//	fmt.Println(issue.Number)
	//}
}

func TestQueryPRInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	prs, _, err := QueryPRInfoByRepo(context.Background(), "cloudwego", "hertz", "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(len(prs))
	for _, pr := range prs {
		fmt.Println(pr)
	}
}

func TestQuerySinglePR(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	pr, err := QuerySinglePR(context.Background(), "PR_kwDOHUxKus44ySzE")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pr)
}

func TestQueryUserInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	res, err := QuerySingleUser(context.Background(), "MDQ6VXNlcjg3NzYwMzM4")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestFor(t *testing.T) {
	for i := range 5 {
		if i == 2 {
			continue
		}
		fmt.Println(i)
	}
	sli := make([]int, 3)
	sli = append(sli, 1, 2)
	fmt.Println(sli) // [0 0 0 1 2]
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().UTC())
	fmt.Println(time.Time{}.IsZero())
}
