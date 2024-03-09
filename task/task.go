package task

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/google/go-github/v60/github"
	"log/slog"
)

// TODO: 记录每个小 task 的执行耗时

func Start() {
	errC := make(chan error)
	// TODO

	if err := util.WaitSignal(errC); err != nil {
		slog.Error("receive close signal error", "signal", err.Error())
		return
	}
	// TODO:
}

func InitTask() {
	// handle groups
	for _, group := range config.GlobalConfig.Groups {
		var (
			groupIssueCount       int
			groupPullRequestCount int
			groupStarCount        int
			groupForkCount        int
			groupContributorCount int
		)
		// handle orgs
		for _, login := range group.Orgs {
			// org data
			org, err := graphql.QueryOrgInfo(context.Background(), login)
			if err != nil {
				slog.Error("error query org info", "err", err.Error())
				continue
			}
			repos, err := graphql.QueryRepoNameByOrg(context.Background(), login)
			if err != nil {
				slog.Error("error query repo name by org", "err", err.Error())
				continue
			}
			var (
				orgIssueCount       int
				orgPullRequestCount int
				orgStarCount        int
				orgForkCount        int
				orgContributorCount int
			)
			// handle repos in org
			for _, nameWithOwner := range repos {
				res, err := InitRepoTask(nameWithOwner)
				if err != nil {
					slog.Error("error execute init repo task", "err", err)
					continue
				}
				// org counting data
				{
					orgIssueCount += res.repo.Issues.TotalCount
					orgPullRequestCount += res.repo.PullRequests.TotalCount
					orgStarCount += res.repo.Stargazers.TotalCount
					orgForkCount += res.repo.Forks.TotalCount
					orgContributorCount += res.contributorCount
				}
				// TODO: insert db repositories, issues, pullrequests, contributors, cursors
			}
			// org in group counting data
			{
				groupIssueCount += orgIssueCount
				groupPullRequestCount += orgPullRequestCount
				groupStarCount += orgStarCount
				groupForkCount += orgForkCount
				groupContributorCount += orgContributorCount
			}
			// TODO: inset db organizations
		}
		// handle repos in group
		for _, nameWithOwner := range group.Repos {
			res, err := InitRepoTask(nameWithOwner)
			if err != nil {
				slog.Error("error execute init repo task", "err", err)
				continue
			}
			// repo in group counting data
			{
				groupIssueCount += res.repo.Issues.TotalCount
				groupPullRequestCount += res.repo.PullRequests.TotalCount
				groupStarCount += res.repo.Stargazers.TotalCount
				groupForkCount += res.repo.Forks.TotalCount
				groupContributorCount += res.contributorCount
			}
			// TODO: insert db repositories
		}
		// TODO: insert db groups
	}
}

func UpdateTask() {
}

// InitRepoTask fetch repo data
func InitRepoTask(nameWithOwner string) (*struct {
	repo             graphql.Repo
	issues           []graphql.Issue
	issueEndCursor   string
	prs              []graphql.PR
	prEndCursor      string
	contributors     []*github.Contributor
	contributorCount int
}, error) {
	owner, name := util.SplitNameWithOwner(nameWithOwner)
	// TODO: use goroutine to optimize
	// repo data
	repo, err := graphql.QueryRepoInfo(context.Background(), owner, name)
	if err != nil {
		return nil, err
	}
	// repo issue data
	issues, issueEndCursor, err := graphql.QueryIssueInfo(context.Background(), owner, name, "")
	if err != nil {
		return nil, err
	}
	// repo pr data
	prs, prEndCursor, err := graphql.QueryPRInfo(context.Background(), owner, name, "")
	if err != nil {
		return nil, err
	}
	// repo contributor data
	contributors, contributorCount, err := rest.GetContributorsByRepo(context.Background(), owner, name)
	if err != nil {
		return nil, err
	}
	return &struct {
		repo             graphql.Repo
		issues           []graphql.Issue
		issueEndCursor   string
		prs              []graphql.PR
		prEndCursor      string
		contributors     []*github.Contributor
		contributorCount int
	}{
		repo:             repo,
		issues:           issues,
		issueEndCursor:   issueEndCursor,
		prs:              prs,
		prEndCursor:      prEndCursor,
		contributors:     contributors,
		contributorCount: contributorCount,
	}, nil
}
