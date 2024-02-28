package db

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"testing"
)

func TestCreateRepo(t *testing.T) {
	Init()
	rows, err := CreateRepository(&model.Repository{
		Owner:            "cloudwego",
		Name:             "hertz",
		IssueCount:       100,
		PullRequestCount: 200,
		StarCount:        300,
		ForkCount:        400,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows)
}
