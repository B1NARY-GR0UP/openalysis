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

// TODO: record execute time of each task
// TODO: add progress bar
// TODO: optimize error handling add retry strategy
// TODO: data cleaning e.g. ByteDance, bytedance, Bytedance => bytedance
// TODO: use transaction

// TODO: 定时任务如果失败（有一个 error 就判定为失败），回退事务，整体重试，通过 chan 来传递和监听是否有 err
// TODO: 使用事务可能需要把 global DB 变为参数传递

func Start(ctx context.Context) {
	errC := make(chan error)

	InitTask(ctx)
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

func Restart(ctx context.Context) {
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
}

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
				}
			}
			// TODO: both success or failed
			contributorCount, err := storage.QueryContributorCountByOrg(ctx, org.ID)
			if err != nil {
				slog.Error("error query contributor count by org", "err", err.Error())
				continue
			}
			if err := storage.CreateOrganization(ctx, &model.Organization{
				Login:            org.Login,
				NodeID:           org.ID,
				IssueCount:       orgCount.IssueCount,
				PullRequestCount: orgCount.PullRequestCount,
				StarCount:        orgCount.StarCount,
				ForkCount:        orgCount.ForkCount,
				ContributorCount: contributorCount,
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
			}
		}
		// TODO: insert groups first, then update counts
		contributorCount, err := storage.QueryContributorCountByGroup(ctx, group.Name)
		if err != nil {
			slog.Error("error query contributor count by group", "err", err.Error())
		}
		if err := storage.CreateGroup(ctx, &model.Group{
			Name:             group.Name,
			IssueCount:       groupCount.IssueCount,
			PullRequestCount: groupCount.PullRequestCount,
			StarCount:        groupCount.StarCount,
			ForkCount:        groupCount.ForkCount,
			ContributorCount: contributorCount,
		}); err != nil {
			slog.Error("error create group", "err", err.Error())
			continue
		}
	}
}

// UpdateTask TODO: fix bug
func UpdateTask(ctx context.Context) {
	for _, group := range config.GlobalConfig.Groups {
		var groupCount Count
		for _, login := range group.Orgs {
			var orgCount Count
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

			_, deleteNeeded := util.CompareSlices(cache[login], repos)
			if err := DeleteRepos(ctx, deleteNeeded); err != nil {
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
				if err := UpdateRepoData(ctx, rd); err != nil {
					slog.Error("error update repo data", "err", err.Error())
					continue
				}
				{
					orgCount.IssueCount += rd.Repo.Issues.TotalCount
					orgCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
					orgCount.StarCount += rd.Repo.Stargazers.TotalCount
					orgCount.ForkCount += rd.Repo.Forks.TotalCount
				}
			}
			contributorCount, err := storage.QueryContributorCountByOrg(ctx, org.ID)
			if err != nil {
				slog.Error("error query contributor count by org", "err", err.Error())
				continue
			}
			if err := storage.UpdateOrganization(ctx, &model.Organization{
				NodeID:           org.ID,
				IssueCount:       orgCount.IssueCount,
				PullRequestCount: orgCount.PullRequestCount,
				StarCount:        orgCount.StarCount,
				ForkCount:        orgCount.ForkCount,
				ContributorCount: contributorCount,
			}); err != nil {
				slog.Error("error update org", "err", err.Error())
				continue
			}
			{
				groupCount.IssueCount += orgCount.IssueCount
				groupCount.PullRequestCount += orgCount.PullRequestCount
				groupCount.StarCount += orgCount.StarCount
				groupCount.ForkCount += orgCount.ForkCount
			}
		}
		for _, nameWithOwner := range group.Repos {
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
			if err := UpdateRepoData(ctx, rd); err != nil {
				slog.Error("error update repo data", "err", err.Error())
				continue
			}
			{
				groupCount.IssueCount += rd.Repo.Issues.TotalCount
				groupCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
				groupCount.StarCount += rd.Repo.Stargazers.TotalCount
				groupCount.ForkCount += rd.Repo.Forks.TotalCount
			}
		}
		contributorCount, err := storage.QueryContributorCountByGroup(ctx, group.Name)
		if err != nil {
			slog.Error("error query contributor count by group", "err", err.Error())
		}
		if err := storage.UpdateGroup(ctx, &model.Group{
			Name:             group.Name,
			IssueCount:       groupCount.IssueCount,
			PullRequestCount: groupCount.PullRequestCount,
			StarCount:        groupCount.StarCount,
			ForkCount:        groupCount.ForkCount,
			ContributorCount: contributorCount,
		}); err != nil {
			slog.Error("error update group", "err", err.Error())
			continue
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
	if err := g.Wait(); err != nil {
		return err
	}
	contributors, contributorCount, err := rest.GetContributorsByRepo(ctx, rd.Owner, rd.Name, rd.Repo.ID)
	if err != nil {
		return err
	}
	rd.Contributors = contributors
	rd.ContributorCount = contributorCount
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

func UpdateRepoData(ctx context.Context, rd *RepoData) error {
	// create repo in each update task due to time series graph
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
			// overlay update issues in db
			if err := storage.UpdateIssue(ctx, &model.Issue{
				NodeID:        issue.ID,
				State:         issue.State,
				IssueClosedAt: issue.ClosedAt,
			}); err != nil {
				return err
			}
		case false:
			// add new issues to db
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
		// assignees of latest issue
		var assignees []*model.IssueAssignees
		for _, assignee := range issue.Assignees.Nodes {
			assignees = append(assignees, &model.IssueAssignees{
				IssueNodeID:    issue.ID,
				IssueNumber:    issue.Number,
				IssueURL:       issue.URL,
				IssueRepoName:  issue.Repository.NameWithOwner,
				AssigneeNodeID: assignee.ID,
				AssigneeLogin:  assignee.Login,
			})
		}
		// handle update in issue_assignees table
		exist, err = storage.IssueAssigneesExist(ctx, issue.ID)
		if err != nil {
			return err
		}
		switch exist {
		// old issues
		case true:
			switch githubv4.IssueState(issue.State) {
			// after update the issue is still open
			case githubv4.IssueStateOpen:
				if util.IsEmptySlice(assignees) {
					// remove from issue_assignees because no assignees
					if err := storage.DeleteIssueAssigneesByIssue(ctx, issue.ID); err != nil {
						return err
					}
				} else {
					// update db if the assignees are changed
					if err := storage.UpdateIssueAssignees(ctx, issue.ID, assignees); err != nil {
						return err
					}
				}
			// after update the issue is closed
			case githubv4.IssueStateClosed:
				// remove from issue_assignees because of closed issue
				if err := storage.DeleteIssueAssigneesByIssue(ctx, issue.ID); err != nil {
					return err
				}
			}
		// new issues
		case false:
			// judge if issue has assignees
			if !util.IsEmptySlice(issue.Assignees.Nodes) && githubv4.IssueState(issue.State) == githubv4.IssueStateOpen {
				// insert into issue_assignees
				if err := storage.CreateIssueAssignees(ctx, assignees); err != nil {
					return err
				}
			}
		}
	}
	// handle pr
	// update old pull requests in db
	// only open pr need to update
	openPRs, err := storage.QueryOPENPullRequests(ctx, rd.Repo.ID)
	if err != nil {
		return err
	}
	for _, openPR := range openPRs {
		// get latest state of old open prs
		pr, err := graphql.QuerySinglePR(ctx, openPR.NodeID)
		if err != nil {
			return err
		}
		// overlay update old open prs
		if err := storage.UpdatePullRequest(ctx, &model.PullRequest{
			NodeID:     pr.ID,
			State:      pr.State,
			PRMergedAt: pr.MergedAt,
			PRClosedAt: pr.ClosedAt,
		}); err != nil {
			return err
		}
		// latest assignees of each old open pr
		var assignees []*model.PullRequestAssignees
		for _, assignee := range pr.Assignees.Nodes {
			assignees = append(assignees, &model.PullRequestAssignees{
				PullRequestNodeID:   pr.ID,
				PullRequestNumber:   pr.Number,
				PullRequestURL:      pr.URL,
				PullRequestRepoName: pr.Repository.NameWithOwner,
				AssigneeNodeID:      assignee.ID,
				AssigneeLogin:       assignee.Login,
			})
		}
		// judge if old pr has assignees
		// NOTE: openPR.NodeID == pr.ID
		exist, err := storage.PullRequestAssigneesExist(ctx, pr.ID)
		if err != nil {
			return err
		}
		switch exist {
		// old open pr has assignees
		case true:
			switch githubv4.PullRequestState(pr.State) {
			// still open
			case githubv4.PullRequestStateOpen:
				if !util.IsEmptySlice(assignees) {
					// if latest pr still have assignees then overlay update
					if err := storage.UpdatePullRequestAssignees(ctx, pr.ID, assignees); err != nil {
						return err
					}
				} else {
					// if latest pr does not have any assignees then remove from pull_request_assignees
					if err := storage.DeletePullRequestAssigneesByPR(ctx, pr.ID); err != nil {
						return err
					}
				}
			// old open pr is closed or merged
			case githubv4.PullRequestStateMerged, githubv4.PullRequestStateClosed:
				if err := storage.DeletePullRequestAssigneesByPR(ctx, pr.ID); err != nil {
					return err
				}
			}
		// old open pr does not have assignees
		case false:
			if !util.IsEmptySlice(assignees) && githubv4.PullRequestState(pr.State) == githubv4.PullRequestStateOpen {
				// latest open pr has assignees then insert into db
				if err := storage.CreatePullRequestAssignees(ctx, assignees); err != nil {
					return err
				}
			}
		}
	}
	// handle new pull requests
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
	if !rd.LastUpdate.IsZero() || rd.EndCursor != "" {
		if err := storage.UpdateCursor(ctx, &model.Cursor{
			LastUpdate: rd.LastUpdate,
			EndCursor:  rd.EndCursor,
		}); err != nil {
			return err
		}
	}
	if err := storage.UpdateOrCreateContributors(ctx, rd.Contributors); err != nil {
		return err
	}
	return nil
}

func DeleteRepos(ctx context.Context, repos []string) error {
	if util.IsEmptySlice(repos) {
		return nil
	}
	for _, repo := range repos {
		owner, name := util.SplitNameWithOwner(repo)
		id, err := storage.QueryRepositoryNodeID(ctx, owner, name)
		if err != nil {
			return err
		}
		if err := storage.DeleteRepository(ctx, id); err != nil {
			return err
		}
		if err := storage.DeleteIssues(ctx, id); err != nil {
			return err
		}
		if err := storage.DeleteIssueAssigneesByRepo(ctx, repo); err != nil {
			return err
		}
		if err := storage.DeletePullRequests(ctx, id); err != nil {
			return err
		}
		if err := storage.DeletePullRequestAssigneesByRepo(ctx, repo); err != nil {
			return err
		}
		if err := storage.DeleteCursor(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
