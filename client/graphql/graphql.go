package graphql

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// TODO: test service procedure

var defaultClient *githubv4.Client

//var query struct {
//	Viewer struct { // must be capital form
//		Login     githubv4.String // use go string type is also ok
//		CreatedAt githubv4.DateTime
//	}
//}

type RepoBasicInfo struct {
	Repository struct {
		Stargazers struct {
			TotalCount int
		}
		Forks struct {
			TotalCount int
		}
		Issues struct {
			TotalCount int
		}
		PullRequests struct {
			TotalCount int
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func GetRepoBasicInfo(ctx context.Context, owner, name string) (*RepoBasicInfo, error) {
	query := &RepoBasicInfo{}
	var variables = map[string]interface{}{
		"owner": githubv4.String(owner), // must use graphql type
		"name":  githubv4.String(name),
	}
	err := defaultClient.Query(ctx, query, variables)
	return query, err
}

// NewClient TODO: use token as argument
func NewClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.GlobalConfig.Backend.Token,
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubv4.NewClient(httpClient)
}
