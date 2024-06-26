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

package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/client/rest"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/storage"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/schollz/progressbar/v3"
	"github.com/shurcooL/githubv4"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Count struct {
	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
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

// TODO: merge InitTask, UpdateTask into DoTask

func InitTask(ctx context.Context, db *gorm.DB) error {
	// handle groups
	for _, group := range config.GlobalConfig.Groups {
		// new progressbar
		bar := progressbar.Default(int64(len(group.Orgs)+len(group.Repos)), "HANDLING GROUP: "+group.Name)

		var groupCount Count
		groupContributions := make(map[string]model.AnalyzedGroupContribution)

		// handle orgs in groups
		for _, login := range group.Orgs {
			// increase bar
			_ = bar.Add(1)

			var orgCount Count
			orgContributions := make(map[string]model.AnalyzedOrgContribution)

			// org data
			org, err := graphql.QueryOrgInfo(ctx, login)
			if err != nil {
				slog.Error("error query org info", "err", err.Error())
				return err
			}
			repos, err := graphql.QueryRepoNameByOrg(ctx, login)
			if err != nil {
				slog.Error("error query repo name by org", "err", err.Error())
				return err
			}

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
					return err
				}
				if err := CreateRepoData(ctx, db, rd); err != nil {
					slog.Error("error create repo data", "err", err.Error())
					return err
				}
				{
					orgCount.IssueCount += rd.Repo.Issues.TotalCount
					orgCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
					orgCount.StarCount += rd.Repo.Stargazers.TotalCount
					orgCount.ForkCount += rd.Repo.Forks.TotalCount
				}
				// handle new contributor
				for _, contributor := range rd.Contributors {
					if _, ok := orgContributions[contributor.NodeID]; !ok {
						orgContributions[contributor.NodeID] = model.AnalyzedOrgContribution{
							Login:         contributor.Login,
							NodeID:        contributor.NodeID,
							OrgLogin:      org.Login,
							OrgNodeID:     org.ID,
							Contributions: 0,
						}
					}
					if _, ok := groupContributions[contributor.NodeID]; !ok {
						groupContributions[contributor.NodeID] = model.AnalyzedGroupContribution{
							Login:         contributor.Login,
							NodeID:        contributor.NodeID,
							GroupName:     group.Name,
							Contributions: 0,
						}
					}
				}
				// calculate contributor contribution by organization and group
				for _, contributor := range rd.Contributors {
					if orgContribution, ok := orgContributions[contributor.NodeID]; ok {
						orgContribution.Contributions += contributor.Contributions
						orgContributions[contributor.NodeID] = orgContribution
					}
					if groupContribution, ok := groupContributions[contributor.NodeID]; ok {
						groupContribution.Contributions += contributor.Contributions
						groupContributions[contributor.NodeID] = groupContribution
					}
				}
			}
			contributorCount, err := storage.QueryContributorCountByOrg(ctx, db, org.ID)
			if err != nil {
				slog.Error("error query contributor count by org", "err", err.Error())
				return err
			}
			if err := storage.CreateOrganization(ctx, db, &model.Organization{
				Login:            org.Login,
				NodeID:           org.ID,
				AvatarURL:        org.AvatarURL,
				IssueCount:       orgCount.IssueCount,
				PullRequestCount: orgCount.PullRequestCount,
				StarCount:        orgCount.StarCount,
				ForkCount:        orgCount.ForkCount,
				ContributorCount: contributorCount,
			}); err != nil {
				slog.Error("error create org", "err", err.Error())
				return err
			}
			if err := storage.CreateGroupsOrganizations(ctx, db, &model.GroupsOrganizations{
				GroupName: group.Name,
				OrgNodeID: org.ID,
			}); err != nil {
				slog.Error("error create group org join", "err", err.Error())
				return err
			}
			{
				groupCount.IssueCount += orgCount.IssueCount
				groupCount.PullRequestCount += orgCount.PullRequestCount
				groupCount.StarCount += orgCount.StarCount
				groupCount.ForkCount += orgCount.ForkCount
			}
			// create analyzed org contributions
			if err := storage.CreateAnalyzedOrgContributions(ctx, db, orgContributions); err != nil {
				slog.Error("error create analyzed org contributions", "err", err.Error())
				return err
			}
		}
		// handle repos in group
		for _, nameWithOwner := range group.Repos {
			// increase bar
			_ = bar.Add(1)

			owner, name := util.SplitNameWithOwner(nameWithOwner)
			rd := &RepoData{
				Owner:         owner,
				Name:          name,
				NameWithOwner: nameWithOwner,
			}
			if err := FetchRepoData(ctx, rd, time.Time{}, ""); err != nil {
				slog.Error("error fetch repo data", "err", err.Error())
				return err
			}
			if err := CreateRepoData(ctx, db, rd); err != nil {
				slog.Error("error create repo data", "err", err.Error())
				return err
			}
			if err := storage.CreateGroupsRepositories(ctx, db, &model.GroupsRepositories{
				GroupName:  group.Name,
				RepoNodeID: rd.Repo.ID,
			}); err != nil {
				slog.Error("error create group repo join", "err", err.Error())
				return err
			}
			{
				groupCount.IssueCount += rd.Repo.Issues.TotalCount
				groupCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
				groupCount.StarCount += rd.Repo.Stargazers.TotalCount
				groupCount.ForkCount += rd.Repo.Forks.TotalCount
			}
			// handle new contributor
			for _, contributor := range rd.Contributors {
				if _, ok := groupContributions[contributor.NodeID]; !ok {
					groupContributions[contributor.NodeID] = model.AnalyzedGroupContribution{
						Login:         contributor.Login,
						NodeID:        contributor.NodeID,
						GroupName:     group.Name,
						Contributions: 0,
					}
				}
			}
			// calculate contributor contribution by group
			for _, contributor := range rd.Contributors {
				if groupContribution, ok := groupContributions[contributor.NodeID]; ok {
					groupContribution.Contributions += contributor.Contributions
					groupContributions[contributor.NodeID] = groupContribution
				}
			}
		}
		contributorCount, err := storage.QueryContributorCountByGroup(ctx, db, group.Name)
		if err != nil {
			slog.Error("error query contributor count by group", "err", err.Error())
			return err
		}
		if err := storage.CreateGroup(ctx, db, &model.Group{
			Name:             group.Name,
			IssueCount:       groupCount.IssueCount,
			PullRequestCount: groupCount.PullRequestCount,
			StarCount:        groupCount.StarCount,
			ForkCount:        groupCount.ForkCount,
			ContributorCount: contributorCount,
		}); err != nil {
			slog.Error("error create group", "err", err.Error())
			return err
		}
		if err := storage.CreateAnalyzedGroupContributions(ctx, db, groupContributions); err != nil {
			slog.Error("err create analyzed group contributions", "err", err.Error())
			return err
		}
	}
	// do cleaner, marker
	if err := CleanContributorCompanyAndLocation(ctx, db); err != nil {
		return err
	}
	// mark must after clean
	if err := MarkContributorCompanyAndLocation(ctx, db); err != nil {
		return err
	}
	return nil
}

func UpdateTask(ctx context.Context, db *gorm.DB) error {
	for _, group := range config.GlobalConfig.Groups {
		// new progressbar
		bar := progressbar.Default(int64(len(group.Orgs)+len(group.Repos)), "HANDLING GROUP: "+group.Name)

		var groupCount Count
		groupContributions := make(map[string]model.AnalyzedGroupContribution)

		for _, login := range group.Orgs {
			// increase bar
			_ = bar.Add(1)

			var orgCount Count
			orgContributions := make(map[string]model.AnalyzedOrgContribution)

			org, err := graphql.QueryOrgInfo(ctx, login)
			if err != nil {
				slog.Error("error query org info", "err", err.Error())
				return err
			}
			repos, err := graphql.QueryRepoNameByOrg(ctx, login)
			if err != nil {
				slog.Error("error query repo name by org", "err", err.Error())
				return err
			}

			oldRepos, err := storage.QueryReposByOrg(ctx, db, org.ID)
			if err != nil {
				slog.Error("error query repos by org", "err", err.Error())
				return err
			}

			_, deleteNeeded := util.CompareSlices(oldRepos, repos)

			// delete repos if org delete it
			if err := DeleteRepos(ctx, db, deleteNeeded); err != nil {
				slog.Error("error delete repos", "err", err.Error())
				return err
			}

			for _, nameWithOwner := range repos {
				owner, name := util.SplitNameWithOwner(nameWithOwner)
				rd := &RepoData{
					Owner:         owner,
					Name:          name,
					NameWithOwner: nameWithOwner,
				}
				cursor, err := storage.QueryCursor(ctx, db, nameWithOwner)
				if err != nil {
					slog.Error("error query cursor", "err", err.Error())
					return err
				}
				if err := FetchRepoData(ctx, rd, cursor.LastUpdate, cursor.EndCursor); err != nil {
					slog.Error("error fetch repo data", "err", err.Error())
					return err
				}
				if err := UpdateRepoData(ctx, db, rd); err != nil {
					slog.Error("error update repo data", "err", err.Error())
					return err
				}
				{
					orgCount.IssueCount += rd.Repo.Issues.TotalCount
					orgCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
					orgCount.StarCount += rd.Repo.Stargazers.TotalCount
					orgCount.ForkCount += rd.Repo.Forks.TotalCount
				}
				// handle new contributor
				for _, contributor := range rd.Contributors {
					if _, ok := orgContributions[contributor.NodeID]; !ok {
						orgContributions[contributor.NodeID] = model.AnalyzedOrgContribution{
							Login:         contributor.Login,
							NodeID:        contributor.NodeID,
							OrgLogin:      org.Login,
							OrgNodeID:     org.ID,
							Contributions: 0,
						}
					}
					if _, ok := groupContributions[contributor.NodeID]; !ok {
						groupContributions[contributor.NodeID] = model.AnalyzedGroupContribution{
							Login:         contributor.Login,
							NodeID:        contributor.NodeID,
							GroupName:     group.Name,
							Contributions: 0,
						}
					}
				}
				// calculate contributor contribution by organization and group
				for _, contributor := range rd.Contributors {
					if orgContribution, ok := orgContributions[contributor.NodeID]; ok {
						orgContribution.Contributions += contributor.Contributions
						orgContributions[contributor.NodeID] = orgContribution
					}
					if groupContribution, ok := groupContributions[contributor.NodeID]; ok {
						groupContribution.Contributions += contributor.Contributions
						groupContributions[contributor.NodeID] = groupContribution
					}
				}
			}
			contributorCount, err := storage.QueryContributorCountByOrg(ctx, db, org.ID)
			if err != nil {
				return err
			}
			if err := storage.CreateOrganization(ctx, db, &model.Organization{
				NodeID:           org.ID,
				Login:            org.Login,
				AvatarURL:        org.AvatarURL,
				IssueCount:       orgCount.IssueCount,
				PullRequestCount: orgCount.PullRequestCount,
				StarCount:        orgCount.StarCount,
				ForkCount:        orgCount.ForkCount,
				ContributorCount: contributorCount,
			}); err != nil {
				slog.Error("error update org", "err", err.Error())
				return err
			}
			{
				groupCount.IssueCount += orgCount.IssueCount
				groupCount.PullRequestCount += orgCount.PullRequestCount
				groupCount.StarCount += orgCount.StarCount
				groupCount.ForkCount += orgCount.ForkCount
			}
			// create analyzed org contributions
			if err := storage.CreateAnalyzedOrgContributions(ctx, db, orgContributions); err != nil {
				slog.Error("error create analyzed org contributions", "err", err.Error())
				return err
			}
		}
		for _, nameWithOwner := range group.Repos {
			// increase bar
			_ = bar.Add(1)

			owner, name := util.SplitNameWithOwner(nameWithOwner)
			rd := &RepoData{
				Owner:         owner,
				Name:          name,
				NameWithOwner: nameWithOwner,
			}
			cursor, err := storage.QueryCursor(ctx, db, nameWithOwner)
			if err != nil {
				slog.Error("error query cursor", "err", err.Error())
				return err
			}
			if err := FetchRepoData(ctx, rd, cursor.LastUpdate, cursor.EndCursor); err != nil {
				slog.Error("error fetch repo data", "err", err.Error())
				return err
			}
			if err := UpdateRepoData(ctx, db, rd); err != nil {
				slog.Error("error update repo data", "err", err.Error())
				return err
			}
			{
				groupCount.IssueCount += rd.Repo.Issues.TotalCount
				groupCount.PullRequestCount += rd.Repo.PullRequests.TotalCount
				groupCount.StarCount += rd.Repo.Stargazers.TotalCount
				groupCount.ForkCount += rd.Repo.Forks.TotalCount
			}
			// handle new contributor
			for _, contributor := range rd.Contributors {
				if _, ok := groupContributions[contributor.NodeID]; !ok {
					groupContributions[contributor.NodeID] = model.AnalyzedGroupContribution{
						Login:         contributor.Login,
						NodeID:        contributor.NodeID,
						GroupName:     group.Name,
						Contributions: 0,
					}
				}
			}
			// calculate contributor contribution by group
			for _, contributor := range rd.Contributors {
				if groupContribution, ok := groupContributions[contributor.NodeID]; ok {
					groupContribution.Contributions += contributor.Contributions
					groupContributions[contributor.NodeID] = groupContribution
				}
			}
		}
		contributorCount, err := storage.QueryContributorCountByGroup(ctx, db, group.Name)
		if err != nil {
			return err
		}
		if err := storage.CreateGroup(ctx, db, &model.Group{
			Name:             group.Name,
			IssueCount:       groupCount.IssueCount,
			PullRequestCount: groupCount.PullRequestCount,
			StarCount:        groupCount.StarCount,
			ForkCount:        groupCount.ForkCount,
			ContributorCount: contributorCount,
		}); err != nil {
			slog.Error("error update group", "err", err.Error())
			return err
		}
		if err := storage.CreateAnalyzedGroupContributions(ctx, db, groupContributions); err != nil {
			slog.Error("err create analyzed group contributions", "err", err.Error())
			return err
		}
	}
	// do cleaner, marker
	if err := CleanContributorCompanyAndLocation(ctx, db); err != nil {
		return err
	}
	// mark must after clean
	if err := MarkContributorCompanyAndLocation(ctx, db); err != nil {
		return err
	}
	return nil
}

func FetchRepoData(ctx context.Context, rd *RepoData, issueCursor time.Time, prCursor string) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		repo, err := graphql.QueryRepoInfo(ctx, rd.Owner, rd.Name)
		if err == nil {
			rd.Repo = repo
		}
		return err
	})
	g.Go(func() error {
		cursor := time.Time{}
		if !issueCursor.IsZero() {
			cursor = issueCursor
		}
		issues, lastUpdate, err := graphql.QueryIssueInfoByRepo(ctx, rd.Owner, rd.Name, cursor)
		if err == nil {
			rd.Issues = issues
			rd.LastUpdate = lastUpdate
		}
		return err
	})
	g.Go(func() error {
		cursor := ""
		if prCursor != "" {
			cursor = prCursor
		}
		prs, endCursor, err := graphql.QueryPRInfoByRepo(ctx, rd.Owner, rd.Name, cursor)
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

func CreateRepoData(ctx context.Context, db *gorm.DB, rd *RepoData) error {
	if err := storage.CreateRepository(ctx, db, &model.Repository{
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
			RepoOwner:      rd.Owner,
			RepoName:       rd.Name,
			Number:         issue.Number,
			State:          issue.State,
			IssueCreatedAt: issue.CreatedAt,
			IssueClosedAt:  util.NilIfZero(issue.ClosedAt),
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
	if err := storage.CreateIssues(ctx, db, issueData); err != nil {
		return err
	}
	if err := storage.CreateIssueAssignees(ctx, db, issueAssignees); err != nil {
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
			RepoOwner:    rd.Owner,
			RepoName:     rd.Name,
			Number:       pr.Number,
			State:        pr.State,
			PRCreatedAt:  pr.CreatedAt,
			PRMergedAt:   util.NilIfZero(pr.MergedAt),
			PRClosedAt:   util.NilIfZero(pr.ClosedAt),
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
	if err := storage.CreatePullRequests(ctx, db, prData); err != nil {
		return err
	}
	if err := storage.CreatePullRequestAssignees(ctx, db, prAssignees); err != nil {
		return err
	}
	if err := storage.CreateCursor(ctx, db, &model.Cursor{
		RepoNodeID:        rd.Repo.ID,
		RepoNameWithOwner: rd.NameWithOwner,
		LastUpdate:        rd.LastUpdate,
		EndCursor:         rd.EndCursor,
	}); err != nil {
		return err
	}
	if err := storage.CreateContributors(ctx, db, rd.Contributors); err != nil {
		return err
	}
	return nil
}

func UpdateRepoData(ctx context.Context, db *gorm.DB, rd *RepoData) error {
	// create repo in each update task due to time series graph
	if err := storage.CreateRepository(ctx, db, &model.Repository{
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
		exist, err := storage.IssueExist(ctx, db, issue.ID)
		if err != nil {
			return err
		}
		switch exist {
		case true:
			// overlay update issues in db
			if err := storage.UpdateIssue(ctx, db, &model.Issue{
				NodeID:        issue.ID,
				State:         issue.State,
				IssueClosedAt: util.NilIfZero(issue.ClosedAt),
			}); err != nil {
				return err
			}
		case false:
			// add new issues to db
			if err := storage.CreateIssues(ctx, db, []*model.Issue{
				{
					NodeID:         issue.ID,
					Author:         issue.Author.Login,
					AuthorNodeID:   issue.Author.User.ID,
					RepoNodeID:     issue.Repository.ID,
					RepoOwner:      rd.Owner,
					RepoName:       rd.Name,
					Number:         issue.Number,
					State:          issue.State,
					IssueCreatedAt: issue.CreatedAt,
					IssueClosedAt:  util.NilIfZero(issue.ClosedAt),
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
		exist, err = storage.IssueAssigneesExist(ctx, db, issue.ID)
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
					if err := storage.DeleteIssueAssigneesByIssue(ctx, db, issue.ID); err != nil {
						return err
					}
				} else {
					// update db if the assignees are changed
					if err := storage.UpdateIssueAssignees(ctx, db, issue.ID, assignees); err != nil {
						return err
					}
				}
			// after update the issue is closed
			case githubv4.IssueStateClosed:
				// remove from issue_assignees because of closed issue
				if err := storage.DeleteIssueAssigneesByIssue(ctx, db, issue.ID); err != nil {
					return err
				}
			}
		// new issues
		case false:
			// judge if issue has assignees
			if !util.IsEmptySlice(issue.Assignees.Nodes) && githubv4.IssueState(issue.State) == githubv4.IssueStateOpen {
				// insert into issue_assignees
				if err := storage.CreateIssueAssignees(ctx, db, assignees); err != nil {
					return err
				}
			}
		}
	}
	// handle pr
	// update old pull requests in db
	// only open pr need to update
	openPRs, err := storage.QueryOPENPullRequests(ctx, db, rd.Repo.ID)
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
		if err := storage.UpdatePullRequest(ctx, db, &model.PullRequest{
			NodeID:     pr.ID,
			State:      pr.State,
			PRMergedAt: util.NilIfZero(pr.MergedAt),
			PRClosedAt: util.NilIfZero(pr.ClosedAt),
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
		exist, err := storage.PullRequestAssigneesExist(ctx, db, pr.ID)
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
					if err := storage.UpdatePullRequestAssignees(ctx, db, pr.ID, assignees); err != nil {
						return err
					}
				} else {
					// if latest pr does not have any assignees then remove from pull_request_assignees
					if err := storage.DeletePullRequestAssigneesByPR(ctx, db, pr.ID); err != nil {
						return err
					}
				}
			// old open pr is closed or merged
			case githubv4.PullRequestStateMerged, githubv4.PullRequestStateClosed:
				if err := storage.DeletePullRequestAssigneesByPR(ctx, db, pr.ID); err != nil {
					return err
				}
			}
		// old open pr does not have assignees
		case false:
			if !util.IsEmptySlice(assignees) && githubv4.PullRequestState(pr.State) == githubv4.PullRequestStateOpen {
				// latest open pr has assignees then insert into db
				if err := storage.CreatePullRequestAssignees(ctx, db, assignees); err != nil {
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
			RepoOwner:    rd.Owner,
			RepoName:     rd.Name,
			Number:       pr.Number,
			State:        pr.State,
			PRCreatedAt:  pr.CreatedAt,
			PRMergedAt:   util.NilIfZero(pr.MergedAt),
			PRClosedAt:   util.NilIfZero(pr.ClosedAt),
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
	if err := storage.CreatePullRequests(ctx, db, prs); err != nil {
		return err
	}
	if err := storage.CreatePullRequestAssignees(ctx, db, prAssignees); err != nil {
		return err
	}
	if err := storage.UpdateOrCreateCursor(ctx, db, &model.Cursor{
		RepoNodeID:        rd.Repo.ID,
		RepoNameWithOwner: rd.NameWithOwner,
		LastUpdate:        rd.LastUpdate,
		EndCursor:         rd.EndCursor,
	}); err != nil {
		return err
	}
	if err := storage.CreateContributors(ctx, db, rd.Contributors); err != nil {
		return err
	}
	return nil
}

func DeleteRepos(ctx context.Context, db *gorm.DB, repos []string) error {
	if util.IsEmptySlice(repos) {
		return nil
	}
	for _, repo := range repos {
		owner, name := util.SplitNameWithOwner(repo)
		id, err := storage.QueryRepositoryNodeID(ctx, db, owner, name)
		if err != nil {
			return err
		}
		if err := storage.DeleteRepository(ctx, db, id); err != nil {
			return err
		}
		if err := storage.DeleteIssues(ctx, db, id); err != nil {
			return err
		}
		if err := storage.DeleteIssueAssigneesByRepo(ctx, db, repo); err != nil {
			return err
		}
		if err := storage.DeletePullRequests(ctx, db, id); err != nil {
			return err
		}
		if err := storage.DeletePullRequestAssigneesByRepo(ctx, db, repo); err != nil {
			return err
		}
		if err := storage.DeleteCursor(ctx, db, id); err != nil {
			return err
		}
	}
	return nil
}

func CleanContributorCompanyAndLocation(ctx context.Context, db *gorm.DB) error {
	if err := storage.UpdateContributorCompanyAndLocation(ctx, db, GlobalCleaner.Clean); err != nil {
		return err
	}
	return nil
}

func MarkContributorCompanyAndLocation(ctx context.Context, db *gorm.DB) error {
	for _, strategy := range GlobalMarker.Marks() {
		if err := storage.UpdateContributorCompanyAndLocationByLogin(ctx, db, strategy.Login, strategy.Company, strategy.Location); err != nil {
			return err
		}
	}
	return nil
}
