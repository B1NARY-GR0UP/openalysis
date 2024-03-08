package graphql

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"time"
)

var GlobalV4Client *githubv4.Client

// Init githubv4 graphql client
func Init() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.GlobalConfig.Backend.Token,
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	GlobalV4Client = githubv4.NewClient(httpClient)
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
	names := make([]string, 0)
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
		names = append(names, node.NameWithOwner)
	}
	return names, nil
}

type RepoInfo struct {
	Repository Repo `graphql:"repository(owner: $owner, name: $name)"`
}

type Repo struct {
	ID    string
	Owner struct {
		ID string
	}
	Issues struct {
		TotalCount int
	}
	PullRequests struct {
		TotalCount int
	}
	Stargazers struct {
		TotalCount int
	}
	Forks struct {
		TotalCount int
	}
}

// QueryRepoInfo return the repo info based on the provided owner and name
func QueryRepoInfo(ctx context.Context, owner, name string) (Repo, error) {
	query := &RepoInfo{}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return Repo{}, err
	}
	return query.Repository, nil
}

type IssueInfo struct {
	Repository struct {
		Issues struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []Issue
		} `graphql:"issues(first: $first, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type Issue struct {
	ID     string
	Author struct {
		Login string
		User  struct { // TODO: handle other types (e.g. bot)
			ID string
		} `graphql:"... on User"`
	}
	Repository struct {
		ID string
	}
	Number    int
	State     string
	CreatedAt time.Time
	ClosedAt  time.Time
}

// QueryIssueInfo return issues according to the repo if endCursor is empty
// it will return the updated issue since last update if endCursor is provided
func QueryIssueInfo(ctx context.Context, owner, name, endCursor string) ([]Issue, string, error) {
	query := &IssueInfo{}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"first": githubv4.Int(100),
		"after": (*githubv4.String)(nil),
	}
	if endCursor != "" {
		variables["after"] = githubv4.NewString(githubv4.String(endCursor))
	}
	issues := make([]Issue, 0)
	for {
		if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
			return nil, "", err
		}
		issues = append(issues, query.Repository.Issues.Nodes...)
		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(githubv4.String(query.Repository.Issues.PageInfo.EndCursor))
	}
	return issues, query.Repository.Issues.PageInfo.EndCursor, nil
}

type PRInfo struct {
	Repository struct {
		PullRequests struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []PR
		} `graphql:"pullRequests(first: $first, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type PR struct {
	ID     string
	Author struct {
		Login string
		User  struct { // TODO: handle other types (e.g. bot)
			ID string
		} `graphql:"... on User"`
	}
	Repository struct {
		ID string
	}
	Number    int
	State     string
	CreatedAt time.Time
	MergedAt  time.Time
	ClosedAt  time.Time
}

// QueryPRInfo return pull requests according to the repo if endCursor is empty
// it will return the updated pull requests since last update if endCursor is provided
func QueryPRInfo(ctx context.Context, owner, name, endCursor string) ([]PR, string, error) {
	query := &PRInfo{}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"first": githubv4.Int(100),
		"after": (*githubv4.String)(nil),
	}
	if endCursor != "" {
		variables["after"] = githubv4.NewString(githubv4.String(endCursor))
	}
	prs := make([]PR, 0)
	for {
		if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
			return nil, "", err
		}
		prs = append(prs, query.Repository.PullRequests.Nodes...)
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(githubv4.String(query.Repository.PullRequests.PageInfo.EndCursor))
	}
	return prs, query.Repository.PullRequests.PageInfo.EndCursor, nil
}
