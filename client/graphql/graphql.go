package graphql

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var GlobalV4Client *githubv4.Client

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
	variables := map[string]interface{}{
		"owner": githubv4.String(owner), // must use graphql type
		"name":  githubv4.String(name),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return nil, err
	}
	return query, nil
}

type RepoName struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []struct {
				NameWithOwner string
			}
		} `graphql:"repositories(first: $first, after: $after)"`
	} `graphql:"organization(login: $login)"`
}

// QueryRepoNameByOrg return repos of the provided org in `org/repo` format
func QueryRepoNameByOrg(ctx context.Context, login string) ([]string, error) {
	query := &RepoName{}
	variables := map[string]interface{}{
		"login": githubv4.String(login),
		"first": githubv4.Int(100),
		"after": (*githubv4.String)(nil),
	}
	nodes := make([]struct{ NameWithOwner string }, 0)
	res := make([]string, 0)
	for {
		if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
			return nil, err
		}
		nodes = append(nodes, query.Organization.Repositories.Nodes...)
		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(githubv4.String(query.Organization.Repositories.PageInfo.EndCursor))
	}
	for _, node := range nodes {
		res = append(res, node.NameWithOwner)
	}
	return res, nil
}

func Init() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.GlobalConfig.Backend.Token,
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	GlobalV4Client = githubv4.NewClient(httpClient)
}
