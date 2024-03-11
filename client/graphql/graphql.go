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
	var (
		repos []struct {
			NameWithOwner string
		}
		names []string
	)
	for {
		if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
			return nil, err
		}
		repos = append(repos, query.Organization.Repositories.Nodes...)
		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(githubv4.String(query.Organization.Repositories.PageInfo.EndCursor))
	}
	for _, repo := range repos {
		names = append(names, repo.NameWithOwner)
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

type OrgInfo struct {
	Organization Org `graphql:"organization(login: $login)"`
}

type Org struct {
	ID    string
	Login string
}

func QueryOrgInfo(ctx context.Context, login string) (Org, error) {
	query := &OrgInfo{}
	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return Org{}, err
	}
	return query.Organization, nil
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
	var issues []Issue
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
	var prs []PR
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

type UserInfo struct {
	Node struct {
		User User `graphql:"... on User"`
	} `graphql:"node(id: $id)"`
}

type User struct {
	Company  string
	Location string
}

func QueryUserInfo(ctx context.Context, nodeID string) (User, error) {
	query := &UserInfo{}
	variables := map[string]interface{}{
		"id": githubv4.ID(nodeID),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return User{}, err
	}
	return query.Node.User, nil
}

// TODO: NOTE: 使用 filterBy 的 since 参数的时候，可以查询到这个时间点后所有有更新的 issues，包括新创建的 issues 和之前状态发送变化的 issues
// TODO: 取消使用 endCursor 来增量更新 issues，使用 since 来更新 issues，但是仍然需要使用分页防止一次获取的数量大于 100 条

// TODO: INIT: query all the open issue with assignees
// TODO: use filterBy: {since: $since, states: $states, assignee: $assignee}
// TODO: $states = "OPEN" $assignee = "*" $since = null

// TODO: UPDATE: query all the updated issue since last updatedAt
// TODO:
