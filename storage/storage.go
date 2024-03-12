package storage

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// TODO: replace dsn with config
	dsn := "root:114514@tcp(localhost:3306)/openalysis?charset=utf8&parseTime=True&loc=Local"
	var err error
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
	if issues == nil || len(issues) == 0 {
		return nil
	}
	return DB.WithContext(ctx).Create(issues).Error
}

func CreatePullRequests(ctx context.Context, prs []*model.PullRequest) error {
	if prs == nil || len(prs) == 0 {
		return nil
	}
	return DB.WithContext(ctx).Create(prs).Error
}

func CreateIssueAssignees(ctx context.Context, as *model.IssueAssignees) error {
	return DB.WithContext(ctx).Create(as).Error
}

func CreatePullRequestAssignees(ctx context.Context, as *model.PullRequestAssignees) error {
	return DB.WithContext(ctx).Create(as).Error
}

func CreateContributors(ctx context.Context, cs []*model.Contributor) error {
	return DB.WithContext(ctx).Create(cs).Error
}

func CreateCursor(ctx context.Context, cursor *model.Cursor) error {
	return DB.WithContext(ctx).Create(cursor).Error
}