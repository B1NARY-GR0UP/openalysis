package model

import (
	"gorm.io/gorm"
	"time"
)

// TODO: 如何在每次更新的时候更新之前的 Issue 的状态
// TODO: 使用增量更新(获取 updatedAt 在上次更新时间之后所有 issues 并更新数据库)或者 webhooks

const (
	IssueStateTypeOpen   = "OPEN"
	IssueStateTypeClosed = "CLOSED"
)

// Issue many to many Contributor
type Issue struct {
	gorm.Model

	NodeID string

	Author       string
	AuthorNodeID string
	RepoNodeID   string
	Number       int

	// OPEN | CLOSED
	State string

	IssueCreatedAt time.Time
	IssueClosedAt  time.Time
}

// TODO: 初始化时查询每个 repo 状态为 OPEN 的 issue，然后获取每个 issue 的 assignees 并插入表
// TODO: 每次更新使用 updatedAt 和 endCursor 获取新增的 issues 和更新的 issues，
// TODO: 从新增的 issues 中选出 OPEN 并且有 assignee 的 issue，使用更新的 issues 更新数据库表的 issues (e.g. delete CLOSED issues)

// IssueAssignees an issue can have multi assignees
// a user can be assigned to multi issues
type IssueAssignees struct {
	gorm.Model
	IssueNodeID    string
	IssueNumber    int
	IssueURL       string
	AssigneeNodeID string
	AssigneeLogin  string
}
