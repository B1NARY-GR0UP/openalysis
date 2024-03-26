package storage

import (
	"context"
	"errors"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := util.AssembleDSN(
		config.GlobalConfig.DataSource.MySQL.Host,
		config.GlobalConfig.DataSource.MySQL.Port,
		config.GlobalConfig.DataSource.MySQL.User,
		config.GlobalConfig.DataSource.MySQL.Password,
		config.GlobalConfig.DataSource.MySQL.Database,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// TODO: use mount
	err = DB.AutoMigrate(
		&model.Cursor{},
		&model.Contributor{},
		&model.Group{},
		&model.Issue{},
		&model.Organization{},
		&model.PullRequest{},
		&model.Repository{},
		&model.GroupsOrganizations{},
		&model.GroupsRepositories{},
		&model.IssueAssignees{},
		&model.PullRequestAssignees{},
	)
	if err != nil {
		panic("failed to migrate tables")
	}
}

func CreateGroup(ctx context.Context, group *model.Group) error {
	return DB.WithContext(ctx).Create(group).Error
}

func CreateOrganization(ctx context.Context, org *model.Organization) error {
	return DB.WithContext(ctx).Create(org).Error
}

func CreateRepository(ctx context.Context, repo *model.Repository) error {
	return DB.WithContext(ctx).Create(repo).Error
}

func CreateGroupsOrganizations(ctx context.Context, join *model.GroupsOrganizations) error {
	return DB.WithContext(ctx).Create(join).Error
}

func CreateGroupsRepositories(ctx context.Context, join *model.GroupsRepositories) error {
	return DB.WithContext(ctx).Create(join).Error
}

func CreateIssues(ctx context.Context, issues []*model.Issue) error {
	if util.IsEmptySlice(issues) {
		return nil
	}
	return DB.WithContext(ctx).Create(issues).Error
}

func IssueExist(ctx context.Context, nodeID string) (bool, error) {
	var issue model.Issue
	if err := DB.WithContext(ctx).Where("node_id = ?", nodeID).First(&issue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateIssue(ctx context.Context, issue *model.Issue) error {
	var currentIssue model.Issue
	if err := DB.WithContext(ctx).Where("node_id = ?", issue.NodeID).First(&currentIssue).Error; err != nil {
		return err
	}
	currentIssue.State = issue.State
	currentIssue.IssueClosedAt = issue.IssueClosedAt
	if err := DB.WithContext(ctx).Save(&currentIssue).Error; err != nil {
		return err
	}
	return nil
}

func CreatePullRequests(ctx context.Context, prs []*model.PullRequest) error {
	if util.IsEmptySlice(prs) {
		return nil
	}
	return DB.WithContext(ctx).Create(prs).Error
}

func UpdatePullRequest(ctx context.Context, pr *model.PullRequest) error {
	var currentPR model.PullRequest
	if err := DB.WithContext(ctx).Where("node_id = ?", pr.NodeID).First(&currentPR).Error; err != nil {
		return err
	}
	currentPR.State = pr.State
	currentPR.PRMergedAt = pr.PRMergedAt
	currentPR.PRClosedAt = pr.PRClosedAt
	if err := DB.WithContext(ctx).Save(&currentPR).Error; err != nil {
		return err
	}
	return nil
}

func QueryOPENPullRequests(ctx context.Context) ([]model.PullRequest, error) {
	var prs []model.PullRequest
	err := DB.WithContext(ctx).Where("state = ?", "OPEN").Find(&prs).Error
	return prs, err
}

func CreateIssueAssignees(ctx context.Context, assignees []*model.IssueAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return DB.WithContext(ctx).Create(assignees).Error
}

func IssueAssigneesExist(ctx context.Context, nodeID string) (bool, error) {
	var assignees model.IssueAssignees
	if err := DB.WithContext(ctx).Where("node_id = ?", nodeID).First(&assignees).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateIssueAssignees(ctx context.Context, issueNodeID string, assignees []model.IssueAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	var currentAssignees []model.IssueAssignees
	if err := DB.WithContext(ctx).Where("issue_node_id = ?", issueNodeID).Find(&currentAssignees).Error; err != nil {
		return err
	}
	more, less := util.CompareSlices(currentAssignees, assignees)
	if err := DB.WithContext(ctx).Create(more).Error; err != nil {
		return err
	}
	for _, e := range less {
		if err := DB.WithContext(ctx).Where("id = ?", e.ID).Delete(&model.IssueAssignees{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteIssueAssignees(ctx context.Context, issueNodeID string) error {
	return DB.WithContext(ctx).Where("issue_node_id = ?", issueNodeID).Delete(&model.IssueAssignees{}).Error
}

func CreatePullRequestAssignees(ctx context.Context, assignees []*model.PullRequestAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return DB.WithContext(ctx).Create(assignees).Error
}

func DeletePullRequestAssignees(ctx context.Context, prNodeID string) error {
	return DB.WithContext(ctx).Where("pull_request_node_id = ?", prNodeID).Delete(&model.PullRequestAssignees{}).Error
}

func CreateContributors(ctx context.Context, cs []*model.Contributor) error {
	if util.IsEmptySlice(cs) {
		return nil
	}
	return DB.WithContext(ctx).Create(cs).Error
}

func CreateCursor(ctx context.Context, cursor *model.Cursor) error {
	return DB.WithContext(ctx).Create(cursor).Error
}

func QueryCursor(ctx context.Context, repo string) (*model.Cursor, error) {
	cursor := &model.Cursor{}
	err := DB.WithContext(ctx).Where("repo_name_with_owner = ?", repo).First(cursor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cursor, nil
	}
	return cursor, err
}
