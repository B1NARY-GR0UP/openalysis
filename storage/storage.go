package storage

import (
	"context"
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

// TODO: handle join table and assignees

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

func CreatePullRequests(ctx context.Context, prs []*model.PullRequest) error {
	if util.IsEmptySlice(prs) {
		return nil
	}
	return DB.WithContext(ctx).Create(prs).Error
}

func CreateIssueAssignees(ctx context.Context, assignees []*model.IssueAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return DB.WithContext(ctx).Create(assignees).Error
}

func CreatePullRequestAssignees(ctx context.Context, assignees []*model.PullRequestAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return DB.WithContext(ctx).Create(assignees).Error
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
	return cursor, err
}
