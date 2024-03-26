package cron

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/robfig/cron/v3"
	"github.com/shurcooL/githubv4"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

// TODO: 记录每个小 task 的执行耗时
// TODO: optimize error handling

func Start(ctx context.Context) {
	errC := make(chan error)

	c := cron.New()
	if _, err := c.AddFunc(config.GlobalConfig.Backend.Cron, func() {
		UpdateTask(ctx)
	}); err != nil {
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

// map[orgLogin][]repoNameWithOwner
var cache map[string][]string

type Count struct {
	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
	ContributorCount int
}

// InitTask TODO: 添加进度条显示
func InitTask(ctx context.Context) {
	// init cache
	cache = make(map[string][]string)
	// handle groups
	for _, group := range config.GlobalConfig.Groups {
		var groupCount Count
		// handle orgs in groups
		for _, login := range group.Orgs {
			var orgCount Count
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

			cache[login] = repos

			// handle repos in org
			// TODO: use errgroup to optimize performance
			for _, nameWithOwner := range repos {
				owner, name := util.SplitNameWithOwner(nameWithOwner)
				rd := &RepoData{
					Owner:         owner,
					Name:          name,
					NameWithOwner: nameWithOwner,
				}
				if err := FetchRepoData(ctx, rd, time.Time{}, ""); err != nil {
					slog.Error("error fetch repo data", "err", err.Error())
					continue
				}
				if err := CreateRepoData(ctx, rd); err != nil {
					slog.Error("error create repo data", "err", err.Error())
					continue
				}
				{
					orgCount.IssueCount += rd.Repo.Issues.TotalCount
					orgCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
					orgCount.StarCount += rd.Repo.Stargazers.TotalCount
					orgCount.ForkCount += rd.Repo.Forks.TotalCount
					orgCount.ContributorCount += rd.ContributorCount
				}
			}
			// TODO: both success or failed
			if err := storage.CreateOrganization(ctx, &model.Organization{
				Login:            org.Login,
				NodeID:           org.ID,
				IssueCount:       orgCount.IssueCount,
				PullRequestCount: orgCount.PullRequestCount,
				StarCount:        orgCount.StarCount,
				ForkCount:        orgCount.ForkCount,
				ContributorCount: orgCount.ContributorCount,
			}); err != nil {
				slog.Error("error create org", "err", err.Error())
				continue
			}
			if err := storage.CreateGroupsOrganizations(ctx, &model.GroupsOrganizations{
				GroupName: group.Name,
				OrgNodeID: org.ID,
			}); err != nil {
				slog.Error("error create group org join", "err", err.Error())
				continue
			}
			{
				groupCount.IssueCount += orgCount.IssueCount
				groupCount.PullRequestCount += orgCount.PullRequestCount
				groupCount.StarCount += orgCount.StarCount
				groupCount.ForkCount += orgCount.ForkCount
				groupCount.ContributorCount += orgCount.ContributorCount
			}
		}
		// handle repos in group
		for _, nameWithOwner := range group.Repos {
			owner, name := util.SplitNameWithOwner(nameWithOwner)
			rd := &RepoData{
				Owner:         owner,
				Name:          name,
				NameWithOwner: nameWithOwner,
			}
			if err := FetchRepoData(ctx, rd, time.Time{}, ""); err != nil {
				slog.Error("error fetch repo data", "err", err.Error())
				continue
			}
			// TODO: both success or failed
			if err := CreateRepoData(ctx, rd); err != nil {
				slog.Error("error create repo data", "err", err.Error())
				continue
			}
			if err := storage.CreateGroupsRepositories(ctx, &model.GroupsRepositories{
				GroupName:  group.Name,
				RepoNodeID: rd.Repo.ID,
			}); err != nil {
				slog.Error("error create group repo join", "err", err.Error())
				continue
			}
			{
				groupCount.IssueCount += rd.Repo.Issues.TotalCount
				groupCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
				groupCount.StarCount += rd.Repo.Stargazers.TotalCount
				groupCount.ForkCount += rd.Repo.Forks.TotalCount
				groupCount.ContributorCount += rd.ContributorCount
			}
		}
		// TODO: insert groups first, then update counts
		if err := storage.CreateGroup(ctx, &model.Group{
			Name:             group.Name,
			IssueCount:       groupCount.IssueCount,
			PullRequestCount: groupCount.PullRequestCount,
			StarCount:        groupCount.StarCount,
			ForkCount:        groupCount.ForkCount,
			ContributorCount: groupCount.ContributorCount,
		}); err != nil {
			slog.Error("error create group", "err", err.Error())
			continue
		}
	}
}

func UpdateTask(ctx context.Context) {
	for _, group := range config.GlobalConfig.Groups {
		var groupCount Count
		for _, login := range group.Orgs {
			var orgCount Count
			// TODO: 处理 org 新增 repo 和删除 repo 的情况
			// TODO: 确保新增的 repo 也进行了处理
			repos, err := graphql.QueryRepoNameByOrg(ctx, login)
			if err != nil {
				slog.Error("error query repo name by org", "err", err.Error())
				continue
			}

			_, deleteNeeded := util.CompareSlices(cache[login], repos)
			if err := DeleteRepos(deleteNeeded); err != nil {
				slog.Error("error delete repos", "err", err.Error())
				continue
			}

			// update cache
			cache[login] = repos

			for _, nameWithOwner := range repos {
				owner, name := util.SplitNameWithOwner(nameWithOwner)
				rd := &RepoData{
					Owner:         owner,
					Name:          name,
					NameWithOwner: nameWithOwner,
				}
				cursor, err := storage.QueryCursor(ctx, nameWithOwner)
				if err != nil {
					slog.Error("error query cursor", "err", err.Error())
					continue
				}
				if err := FetchRepoData(ctx, rd, cursor.LastUpdate, cursor.EndCursor); err != nil {
					slog.Error("error fetch repo data", "err", err.Error())
					continue
				}
				// TODO
				if err := UpdateRepoData(ctx, rd); err != nil {
					slog.Error("error update repo data", "err", err.Error())
					continue
				}
				{
					orgCount.IssueCount += rd.Repo.Issues.TotalCount
					orgCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
					orgCount.StarCount += rd.Repo.Stargazers.TotalCount
					orgCount.ForkCount += rd.Repo.Forks.TotalCount
					orgCount.ContributorCount += rd.ContributorCount
				}
			}
		}
	}
}

type RepoData struct {
	Owner            string
	Name             string
	NameWithOwner    string
	Repo             graphql.Repo
	Issues           []graphql.Issue
	LastUpdate       time.Time
	PRs              []graphql.PR
	EndCursor        string
	Contributors     []*model.Contributor
	ContributorCount int
}

func FetchRepoData(ctx context.Context, rd *RepoData, lu time.Time, ec string) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		repo, err := graphql.QueryRepoInfo(ctx, rd.Owner, rd.Name)
		if err == nil {
			rd.Repo = repo
		}
		return err
	})
	g.Go(func() error {
		t := time.Time{}
		if !lu.IsZero() {
			t = lu
		}
		issues, lastUpdate, err := graphql.QueryIssueInfoByRepo(ctx, rd.Owner, rd.Name, t)
		if err == nil {
			rd.Issues = issues
			rd.LastUpdate = lastUpdate
		}
		return err
	})
	g.Go(func() error {
		c := ""
		if ec != "" {
			c = ec
		}
		prs, endCursor, err := graphql.QueryPRInfoByRepo(ctx, rd.Owner, rd.Name, c)
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
	if err := storage.CreateRepository(ctx, &model.Repository{
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
	var issueAssignees []*model.IssueAssignees
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
		if !util.IsEmptySlice(issue.Assignees.Nodes) && githubv4.IssueState(issue.State) == githubv4.IssueStateOpen {
			for _, assignee := range issue.Assignees.Nodes {
				issueAssignees = append(issueAssignees, &model.IssueAssignees{
					IssueNodeID:    issue.ID,
					IssueNumber:    issue.Number,
					IssueURL:       issue.URL,
					IssueRepoName:  issue.Repository.NameWithOwner,
					AssigneeNodeID: assignee.ID,
					AssigneeLogin:  assignee.Login,
				})
			}
		}
	}
	if err := storage.CreateIssues(ctx, issueData); err != nil {
		return err
	}
	if err := storage.CreateIssueAssignees(ctx, issueAssignees); err != nil {
		return err
	}
	var prData []*model.PullRequest
	var prAssignees []*model.PullRequestAssignees
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
		if !util.IsEmptySlice(pr.Assignees.Nodes) && githubv4.PullRequestState(pr.State) == githubv4.PullRequestStateOpen {
			for _, assignee := range pr.Assignees.Nodes {
				prAssignees = append(prAssignees, &model.PullRequestAssignees{
					PullRequestNodeID:   pr.ID,
					PullRequestNumber:   pr.Number,
					PullRequestURL:      pr.URL,
					PullRequestRepoName: pr.Repository.NameWithOwner,
					AssigneeNodeID:      assignee.ID,
					AssigneeLogin:       assignee.Login,
				})
			}
		}
	}
	if err := storage.CreatePullRequests(ctx, prData); err != nil {
		return err
	}
	if err := storage.CreatePullRequestAssignees(ctx, prAssignees); err != nil {
		return err
	}
	if err := storage.CreateCursor(context.Background(), &model.Cursor{
		RepoNodeID:        rd.Repo.ID,
		RepoNameWithOwner: rd.NameWithOwner,
		LastUpdate:        rd.LastUpdate,
		EndCursor:         rd.EndCursor,
	}); err != nil {
		return err
	}
	if err := storage.CreateContributors(ctx, rd.Contributors); err != nil {
		return err
	}
	return nil
}

// TODO: INIT 全部插入 issues table; 在循环中判断，如果有 assignees 且状态是 OPEN 则插入 assignees table
//       UPDATE 判断 issues table 中是否存在，如果存在则进行覆盖更新，如果不存在则插入 issues table; 在循环中判断是否在 assignees table 中，
//       如果在并且 graphql 查询得到的状态为 OPEN 就覆盖更新，CLOSED 就删除，如果不在则判断是否有 assignees 且状态是 OPEN，都满足的话插入 assignees table

// TODO: INIT 全部插入 prs table; 在循环中判断，如果有 assignees 且状态是 OPEN 则插入 assignees table
//       UPDATE 全部插入 prs table，循环判断是否有 OPEN 并且有 assignees 的 pr，如果有则插入 assignees table;
//       查询 prs table 中 state 为 OPEN 的 prs，查询后覆盖更新数据库，查询 assignees table 中的所有 prs，如果 graphql query
//       后返回的 state 为 CLOSED 或 MERGED，则从 prs table 中删除，否则覆盖更新。

func UpdateRepoData(ctx context.Context, rd *RepoData) error {
	if err := storage.CreateRepository(ctx, &model.Repository{
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
	// handle issue
	for _, issue := range rd.Issues {
		// handle update in issues table
		exist, err := storage.IssueExist(ctx, issue.ID)
		if err != nil {
			return err
		}
		switch exist {
		case true:
			if err := storage.UpdateIssue(ctx, &model.Issue{
				NodeID:        issue.ID,
				State:         issue.State,
				IssueClosedAt: issue.ClosedAt,
			}); err != nil {
				return err
			}
		case false:
			if err := storage.CreateIssues(ctx, []*model.Issue{
				{
					NodeID:         issue.ID,
					Author:         issue.Author.Login,
					AuthorNodeID:   issue.Author.User.ID,
					RepoNodeID:     issue.Repository.ID,
					Number:         issue.Number,
					State:          issue.State,
					IssueCreatedAt: issue.CreatedAt,
					IssueClosedAt:  issue.ClosedAt,
				},
			}); err != nil {
				return err
			}
		}
		// handle update in issue_assignees table
		exist, err = storage.IssueAssigneesExist(ctx, issue.ID)
		if err != nil {
			return err
		}
		switch exist {
		case true:
			switch githubv4.IssueState(issue.State) {
			case githubv4.IssueStateOpen:
				var assignees []model.IssueAssignees
				for _, assignee := range issue.Assignees.Nodes {
					assignees = append(assignees, model.IssueAssignees{
						IssueNodeID:    issue.ID,
						IssueNumber:    issue.Number,
						IssueURL:       issue.URL,
						IssueRepoName:  issue.Repository.NameWithOwner,
						AssigneeNodeID: assignee.ID,
						AssigneeLogin:  assignee.Login,
					})
				}
				if err := storage.UpdateIssueAssignees(ctx, issue.ID, assignees); err != nil {
					return err
				}
			case githubv4.IssueStateClosed:
				if err := storage.DeleteIssueAssignees(ctx, issue.ID); err != nil {
					return err
				}
			}
		case false:
			if !util.IsEmptySlice(issue.Assignees.Nodes) && githubv4.IssueState(issue.State) == githubv4.IssueStateOpen {
				var issueAssignees []*model.IssueAssignees
				for _, assignee := range issue.Assignees.Nodes {
					issueAssignees = append(issueAssignees, &model.IssueAssignees{
						IssueNodeID:    issue.ID,
						IssueNumber:    issue.Number,
						IssueURL:       issue.URL,
						IssueRepoName:  issue.Repository.NameWithOwner,
						AssigneeNodeID: assignee.ID,
						AssigneeLogin:  assignee.Login,
					})
				}
				if err := storage.CreateIssueAssignees(ctx, issueAssignees); err != nil {
					return err
				}
			}
		}
	}
	// handle pr
	var prs []*model.PullRequest
	var prAssignees []*model.PullRequestAssignees
	for _, pr := range rd.PRs {
		prs = append(prs, &model.PullRequest{
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
		// handle update in pull_request_assignees table
		if !util.IsEmptySlice(pr.Assignees.Nodes) && githubv4.PullRequestState(pr.State) == githubv4.PullRequestStateOpen {
			for _, assignee := range pr.Assignees.Nodes {
				prAssignees = append(prAssignees, &model.PullRequestAssignees{
					PullRequestNodeID:   pr.ID,
					PullRequestNumber:   pr.Number,
					PullRequestURL:      pr.URL,
					PullRequestRepoName: pr.Repository.NameWithOwner,
					AssigneeNodeID:      assignee.ID,
					AssigneeLogin:       assignee.Login,
				})
			}
		}
	}
	// handle update in pull_requests table
	if err := storage.CreatePullRequests(ctx, prs); err != nil {
		return err
	}
	if err := storage.CreatePullRequestAssignees(ctx, prAssignees); err != nil {
		return err
	}
	openPRs, err := storage.QueryOPENPullRequests(ctx)
	if err != nil {
		return err
	}
	for _, openPR := range openPRs {
		pr, err := graphql.QuerySinglePR(ctx, openPR.NodeID)
		if err != nil {
			return err
		}
		if err := storage.UpdatePullRequest(ctx, &model.PullRequest{
			NodeID:     pr.ID,
			State:      pr.State,
			PRMergedAt: pr.MergedAt,
			PRClosedAt: pr.ClosedAt,
		}); err != nil {
			return err
		}
		if state := githubv4.PullRequestState(pr.State); state == githubv4.PullRequestStateMerged || state == githubv4.PullRequestStateClosed {
			if err := storage.DeletePullRequestAssignees(ctx, pr.ID); err != nil {
				return err
			}
		}
	}
	if err := storage.UpdateCursor(ctx, &model.Cursor{
		LastUpdate: rd.LastUpdate,
		EndCursor:  rd.EndCursor,
	}); err != nil {
		return err
	}
	if err := storage.UpdateOrCreateContributors(ctx, rd.Contributors); err != nil {
		return err
	}
	return nil
}

func DeleteRepos(repos []string) error {
	// TODO
	return nil
}
