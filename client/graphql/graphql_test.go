package graphql

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
)

func TestSendRequest(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	info, err := GetRepoBasicInfo(context.Background(), "cloudwego", "hertz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

// TODO: fix bug
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
