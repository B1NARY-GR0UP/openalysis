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

// IssueAssignees an issue can have multi assignees
// a user can be assigned to multi issues
type IssueAssignees struct {
	gorm.Model

	IssueNodeID   string
	IssueNumber   int
	IssueURL      string
	IssueRepoName string

	AssigneeNodeID string
	AssigneeLogin  string
}
