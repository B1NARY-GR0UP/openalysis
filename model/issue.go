package model

import (
	"gorm.io/gorm"
	"time"
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
