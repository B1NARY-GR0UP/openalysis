package graphql

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
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
	//issues, endCursor, err := QueryIssueInfo(context.Background(), "rainiring", "test", "")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(len(issues))
	//fmt.Println(endCursor) // use for next update
	//for _, issue := range issues {
	//	fmt.Println(issue.Number)
	//}

	issues, _, err := QueryIssueInfo(context.Background(), "rainiring", "test", "Y3Vyc29yOnYyOpHOgar3bA==")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(len(issues))
	for _, issue := range issues {
		fmt.Println(issue.Number)
	}
}

func TestQueryPRInfo(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	prs, _, err := QueryPRInfo(context.Background(), "cloudwego", "hertz", "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(len(prs))
	for _, pr := range prs {
		fmt.Println(pr)
	}
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
