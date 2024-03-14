package task

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

// TODO: 记录每个小 task 的执行耗时

func Start() {
	errC := make(chan error)

	c := cron.New()
	if _, err := c.AddFunc(config.GlobalConfig.Backend.Frequency, func() {}); err != nil {
		slog.Error("error doing cron", "err", err)
		errC <- err
	}
	c.Start()
	defer c.Stop()

	if err := util.WaitSignal(errC); err != nil {
		slog.Error("receive close signal error", "signal", err.Error())
		return
	}
	slog.Info("openalysis service stopped")
}

// InitTask TODO: 添加进度条显示
func InitTask(ctx context.Context) {
	//groupBar := progressbar.Default(int64(len(config.GlobalConfig.Groups)), "Handle Groups")
	// handle groups
	for _, group := range config.GlobalConfig.Groups {
		//_ = groupBar.Add(1)
		var (
			groupIssueCount       int
			groupPullRequestCount int
			groupStarCount        int
			groupForkCount        int
			groupContributorCount int
		)
		//orgBar := progressbar.Default(int64(len(group.Orgs)), fmt.Sprintf("Handle Org in Group [%v]", group.Name))
		// handle orgs in groups
		for _, login := range group.Orgs {
			//_ = orgBar.Add(1)
			var (
				orgIssueCount       int
				orgPullRequestCount int
				orgStarCount        int
				orgForkCount        int
				orgContributorCount int
			)
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
			// handle repos in org
			// TODO: use errgroup to optimize performance
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
				{
					orgIssueCount += rd.Repo.Issues.TotalCount
					orgPullRequestCount += rd.Repo.PullRequests.TotalCount
					orgStarCount += rd.Repo.Stargazers.TotalCount
					orgForkCount += rd.Repo.Forks.TotalCount
					orgContributorCount += rd.ContributorCount
				}
			}
			// TODO: both success or failed
			if err := storage.CreateOrganization(ctx, &model.Organization{
				Login:            org.Login,
				NodeID:           org.ID,
				IssueCount:       orgIssueCount,
				PullRequestCount: orgPullRequestCount,
				StarCount:        orgStarCount,
				ForkCount:        orgForkCount,
				ContributorCount: orgContributorCount,
			}); err != nil {
				slog.Error("error create org", "err", err)
				continue
			}
			if err := storage.CreateGroupsOrganizations(ctx, &model.GroupsOrganizations{
				GroupName: group.Name,
				OrgNodeID: org.ID,
			}); err != nil {
				slog.Error("error create group org join", "err", err)
				continue
			}
			{
				groupIssueCount += orgIssueCount
				groupPullRequestCount += orgPullRequestCount
				groupStarCount += orgStarCount
				groupForkCount += orgForkCount
				groupContributorCount += orgContributorCount
			}
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
			// TODO: both success or failed
			if err := CreateRepoData(ctx, rd); err != nil {
				slog.Error("error create repo data", "err", err)
				continue
			}
			if err := storage.CreateGroupsRepositories(ctx, &model.GroupsRepositories{
				GroupName:  group.Name,
				RepoNodeID: rd.Repo.ID,
			}); err != nil {
				slog.Error("error create group repo join", "err", err)
				continue
			}
			{
				groupIssueCount += rd.Repo.Issues.TotalCount
				groupPullRequestCount += rd.Repo.PullRequests.TotalCount
				groupStarCount += rd.Repo.Stargazers.TotalCount
				groupForkCount += rd.Repo.Forks.TotalCount
				groupContributorCount += rd.ContributorCount
			}
		}
		// TODO: insert groups first, then update counts
		if err := storage.CreateGroup(ctx, &model.Group{
			Name:             group.Name,
			IssueCount:       groupIssueCount,
			PullRequestCount: groupPullRequestCount,
			StarCount:        groupStarCount,
			ForkCount:        groupForkCount,
			ContributorCount: groupContributorCount,
		}); err != nil {
			slog.Error("error create group", "err", err)
			continue
		}
	}
}

// TODO: update task 是不是可以简化流程，例如在 init task 时提前把所有 repo 存储到内存

func UpdateTask() {
}

type RepoData struct {
	Owner            string
	Name             string
	Repo             graphql.Repo
	Issues           []graphql.Issue
	LastUpdate       time.Time
	PRs              []graphql.PR
	EndCursor        string
	Contributors     []*model.Contributor
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
		issues, lastUpdate, err := graphql.QueryIssueInfoByRepo(ctx, rd.Owner, rd.Name, time.Time{})
		if err == nil {
			rd.Issues = issues
			rd.LastUpdate = lastUpdate
		}
		return err
	})
	g.Go(func() error {
		prs, endCursor, err := graphql.QueryPRInfoByRepo(ctx, rd.Owner, rd.Name, "")
		if err == nil {
			rd.PRs = prs
			rd.EndCursor = endCursor
		}
		return err
	})
	g.Go(func() error {
		contributors, contributorCount, err := rest.GetContributorsByRepo(ctx, rd.Owner, rd.Name, rd.Repo.ID)
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
	if err := storage.CreateRepository(context.Background(), &model.Repository{
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
		return err
	}
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
	if err := storage.CreateIssues(ctx, issueData); err != nil {
		return err
	}
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
	if err := storage.CreatePullRequests(ctx, prData); err != nil {
		return err
	}
	if err := storage.CreateCursor(context.Background(), &model.Cursor{
		RepoNodeID: rd.Repo.ID,
		LastUpdate: rd.LastUpdate,
		Type:       model.CursorTypeIssue,
	}); err != nil {
		return err
	}
	if err := storage.CreateCursor(context.Background(), &model.Cursor{
		RepoNodeID: rd.Repo.ID,
		EndCursor:  rd.EndCursor,
		Type:       model.CursorTypePR,
	}); err != nil {
		return err
	}
	if err := storage.CreateContributors(ctx, rd.Contributors); err != nil {
		return err
	}
	return nil
}
