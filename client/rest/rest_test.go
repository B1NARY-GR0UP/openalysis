package rest

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
)

func TestGetContributors(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	Init()
	res, count, err := GetContributorsByRepo(context.Background(), "cloudwego", "hertz")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(count)
	for _, contributor := range res {
		fmt.Println(contributor.GetLogin())
	}
}
