package graphql

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
)

func TestSendRequest(t *testing.T) {
	if err := config.DefaultConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	defaultClient = NewClient()
	info, err := GetRepoBasicInfo(context.Background(), "cloudwego", "hertz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
