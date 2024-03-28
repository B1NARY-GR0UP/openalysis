package rest

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
)

func TestGetContributors(t *testing.T) {
	if err := config.GlobalConfig.ReadInConfig("../../default.yaml"); err != nil {
		panic(err.Error())
	}
	graphql.Init()
	Init()
	res, count, err := GetContributorsByRepo(context.Background(), "cloudwego", "hertz", "777")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(count)
	for _, contributor := range res {
		fmt.Println(contributor.Contributions)
	}
}
