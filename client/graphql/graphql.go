// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphql

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

// TODO: separate methods into multi files according to the used model
// TODO: optimize graphql (or storage) model, provide more general data model for customized

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
	ID        string
	Login     string
	AvatarURL string
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
		} `graphql:"issues(first: $issuesFirst, after: $issuesAfter, filterBy: { since: $since })"`
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
		ID            string
		NameWithOwner string
	}
	Number    int
	URL       string
	State     string
	CreatedAt time.Time
	ClosedAt  time.Time
	Assignees struct { // TODO: handle paging (we assume that a issue's assignees count is less than 100)
		Nodes []IssueAssignee
	} `graphql:"assignees(first: $assigneesFirst, after: $assigneesAfter)"`
}

type IssueAssignee struct {
	ID    string
	Login string
}

// QueryIssueInfoByRepo return issues according to the repo if lastUpdate is empty
// it will return the issues since last update if lastUpdate is provided
// including new issues and updated issues
func QueryIssueInfoByRepo(ctx context.Context, owner, name string, lastUpdate time.Time) ([]Issue, time.Time, error) {
	query := &IssueInfo{}
	variables := map[string]interface{}{
		"owner":          githubv4.String(owner),
		"name":           githubv4.String(name),
		"issuesFirst":    githubv4.Int(100),
		"issuesAfter":    (*githubv4.String)(nil),
		"since":          (*githubv4.DateTime)(nil),
		"assigneesFirst": githubv4.Int(100),
		"assigneesAfter": (*githubv4.String)(nil),
	}
	if !lastUpdate.IsZero() {
		variables["since"] = githubv4.NewDateTime(githubv4.DateTime{Time: lastUpdate})
	}
	var issues []Issue
	for {
		if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
			return nil, time.Time{}, err
		}
		issues = append(issues, query.Repository.Issues.Nodes...)
		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["issuesAfter"] = githubv4.NewString(githubv4.String(query.Repository.Issues.PageInfo.EndCursor))
	}
	return issues, time.Now().UTC(), nil
}

type PRInfo struct {
	Repository struct {
		PullRequests struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []PR
		} `graphql:"pullRequests(first: $prFirst, after: $prAfter)"`
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
		ID            string
		NameWithOwner string
	}
	Number    int
	URL       string
	State     string
	CreatedAt time.Time
	MergedAt  time.Time
	ClosedAt  time.Time
	Assignees struct { // TODO: handle paging (we assume that a pr's assignees count is less than 100)
		Nodes []PRAssignee
	} `graphql:"assignees(first: $assigneesFirst, after: $assigneesAfter)"`
}

type PRAssignee struct {
	ID    string
	Login string
}

// QueryPRInfoByRepo return pull requests according to the repo if lastUpdate is empty
// it will return the prs since last update if lastUpdate is provided
// including new prs and updated prs
func QueryPRInfoByRepo(ctx context.Context, owner, name, endCursor string) ([]PR, string, error) {
	query := &PRInfo{}
	variables := map[string]interface{}{
		"owner":          githubv4.String(owner),
		"name":           githubv4.String(name),
		"prFirst":        githubv4.Int(100),
		"prAfter":        (*githubv4.String)(nil),
		"assigneesFirst": githubv4.Int(100),
		"assigneesAfter": (*githubv4.String)(nil),
	}
	if endCursor != "" {
		variables["prAfter"] = githubv4.NewString(githubv4.String(endCursor))
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
		variables["prAfter"] = githubv4.NewString(githubv4.String(query.Repository.PullRequests.PageInfo.EndCursor))
	}
	cursor := query.Repository.PullRequests.PageInfo.EndCursor
	if cursor == "" {
		cursor = endCursor
	}
	return prs, cursor, nil
}

type SinglePR struct {
	Node struct {
		PullRequest PR `graphql:"... on PullRequest"`
	} `graphql:"node(id: $id)"`
}

func QuerySinglePR(ctx context.Context, id string) (PR, error) {
	query := &SinglePR{}
	variables := map[string]interface{}{
		"id":             githubv4.ID(id),
		"assigneesFirst": githubv4.Int(100),
		"assigneesAfter": (*githubv4.String)(nil),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return PR{}, err
	}
	return query.Node.PullRequest, nil
}

type SingleUser struct {
	Node struct {
		User User `graphql:"... on User"`
	} `graphql:"node(id: $id)"`
}

type User struct {
	Company  string
	Location string
}

func QuerySingleUser(ctx context.Context, nodeID string) (User, error) {
	query := &SingleUser{}
	variables := map[string]interface{}{
		"id": githubv4.ID(nodeID),
	}
	if err := GlobalV4Client.Query(ctx, query, variables); err != nil {
		return User{}, err
	}
	return query.Node.User, nil
}
