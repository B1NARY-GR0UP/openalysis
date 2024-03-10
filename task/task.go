package task

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/db"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/google/go-github/v60/github"
	"golang.org/x/sync/errgroup"
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

func InitTask(ctx context.Context) {
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
			org, err := graphql.QueryOrgInfo(ctx, login)
			if err != nil {
				slog.Error("error query org info", "err", err.Error())
				continue
			}
			repos, err := graphql.QueryRepoNameByOrg(ctx, login)
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
				owner, name := util.SplitNameWithOwner(nameWithOwner)
				rd := &RepoData{
					Owner: owner,
					Name:  name,
				}
				if err := FetchRepoData(ctx, rd); err != nil {
					slog.Error("error fetch repo data", "err", err)
					continue
				}
				if err := CreateRepoData(ctx, rd); err != nil {
					slog.Error("error create repo data", "err", err)
					continue
				}
				// org counting data
				{
					orgIssueCount += rd.Repo.Issues.TotalCount
					orgPullRequestCount += rd.Repo.PullRequests.TotalCount
					orgStarCount += rd.Repo.Stargazers.TotalCount
					orgForkCount += rd.Repo.Forks.TotalCount
					orgContributorCount += rd.ContributorCount
				}
			}
			// org in group counting data
			{
				groupIssueCount += orgIssueCount
				groupPullRequestCount += orgPullRequestCount
				groupStarCount += orgStarCount
				groupForkCount += orgForkCount
				groupContributorCount += orgContributorCount
			}
			// TODO: insert db organizations
		}
		// handle repos in group
		for _, nameWithOwner := range group.Repos {
			owner, name := util.SplitNameWithOwner(nameWithOwner)
			rd := &RepoData{
				Owner: owner,
				Name:  name,
			}
			if err := FetchRepoData(ctx, rd); err != nil {
				slog.Error("error fetch repo data", "err", err)
				continue
			}
			// repo in group counting data
			{
				groupIssueCount += rd.Repo.Issues.TotalCount
				groupPullRequestCount += rd.Repo.PullRequests.TotalCount
				groupStarCount += rd.Repo.Stargazers.TotalCount
				groupForkCount += rd.Repo.Forks.TotalCount
				groupContributorCount += rd.ContributorCount
			}
			// TODO: insert db repositories
		}
		// TODO: insert db groups
	}
}

func UpdateTask() {
}

type RepoData struct {
	Owner            string
	Name             string
	Repo             graphql.Repo
	Issues           []graphql.Issue
	IssueEndCursor   string
	PRs              []graphql.PR
	PREndCursor      string
	Contributors     []*github.Contributor
	ContributorCount int
}

func FetchRepoData(ctx context.Context, rd *RepoData) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		repo, err := graphql.QueryRepoInfo(ctx, rd.Owner, rd.Name)
		if err == nil {
			rd.Repo = repo
		}
		return err
	})
	g.Go(func() error {
		issues, issueEndCursor, err := graphql.QueryIssueInfo(ctx, rd.Owner, rd.Name, "")
		if err == nil {
			rd.Issues = issues
			rd.IssueEndCursor = issueEndCursor
		}
		return err
	})
	g.Go(func() error {
		prs, prEndCursor, err := graphql.QueryPRInfo(ctx, rd.Owner, rd.Name, "")
		if err == nil {
			rd.PRs = prs
			rd.PREndCursor = prEndCursor
		}
		return err
	})
	g.Go(func() error {
		contributors, contributorCount, err := rest.GetContributorsByRepo(ctx, rd.Owner, rd.Name)
		if err == nil {
			rd.Contributors = contributors
			rd.ContributorCount = contributorCount
		}
		return err
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func CreateRepoData(ctx context.Context, rd *RepoData) error {
	// TODO: insert db repositories, issues, pullrequests, contributors, cursors
	if err := db.CreateRepository(context.Background(), &model.Repository{
		Owner:            rd.Owner,
		Name:             rd.Name,
		NodeID:           rd.Repo.ID,
		OwnerNodeID:      rd.Repo.Owner.ID,
		IssueCount:       rd.Repo.Issues.TotalCount,
		PullRequestCount: rd.Repo.PullRequests.TotalCount,
		StarCount:        rd.Repo.Stargazers.TotalCount,
		ForkCount:        rd.Repo.Forks.TotalCount,
		ContributorCount: rd.ContributorCount,
	}); err != nil {
		slog.Error("error create repo", "err", err)
	}
	db.CreateCursor(context.Background(), &model.Cursor{
		EndCursor: rd.IssueEndCursor,
		Type:      model.CursorTypeIssue,
	})
	db.CreateCursor(context.Background(), &model.Cursor{
		EndCursor: rd.PREndCursor,
		Type:      model.CursorTypePR,
	})
	var issueData []*model.Issue
	for _, issue := range rd.Issues {
		issueData = append(issueData, &model.Issue{
			NodeID:         issue.ID,
			Author:         issue.Author.Login,
			AuthorNodeID:   issue.Author.User.ID,
			RepoNodeID:     issue.Repository.ID,
			Number:         issue.Number,
			State:          issue.State,
			IssueCreatedAt: issue.CreatedAt,
			IssueClosedAt:  issue.ClosedAt,
		})
	}
	db.CreateIssues(ctx, issueData)
	var prData []*model.PullRequest
	for _, pr := range rd.PRs {
		prData = append(prData, &model.PullRequest{
			NodeID:       pr.ID,
			Author:       pr.Author.Login,
			AuthorNodeID: pr.Author.User.ID,
			RepoNodeID:   pr.Repository.ID,
			Number:       pr.Number,
			State:        pr.State,
			PRCreatedAt:  pr.CreatedAt,
			PRMergedAt:   pr.MergedAt,
			PRClosedAt:   pr.ClosedAt,
		})
	}
	db.CreatePullRequests(ctx, prData)
	var contributorsData []*model.Contributor
	// TODO: 包装一个自己的 Contributor 类型，这个类型在 rest 调用 graphql 后返回，其中包含了数据库表中所有需要的数据
	for _, contributor := range rd.Contributors {
		contributorsData = append(contributorsData, &model.Contributor{
			Login:         contributor.GetLogin(),
			NodeID:        contributor.GetNodeID(),
			Company:       "",
			Location:      "",
			AvatarURL:     contributor.GetAvatarURL(),
			RepoNodeID:    rd.Repo.ID,
			Contributions: contributor.GetContributions(),
		})
	}
}
