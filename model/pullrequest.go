package model

import (
	"gorm.io/gorm"
	"time"
)

// PullRequest many to many Contributor
type PullRequest struct {
	gorm.Model

	NodeID string

	Author       string
	AuthorNodeID string
	RepoNodeID   string
	Number       int

	// CLOSED | MERGED | OPEN
	State string

	PRCreatedAt time.Time
	PRMergedAt  time.Time
	PRClosedAt  time.Time
}

// PullRequestAssignees a pr can have multi reviewers
// a user can be assigned to multi pull requests
type PullRequestAssignees struct {
	gorm.Model
	PullRequestNodeID string
	AssigneeNodeID    string
	AssigneeLogin     string
}
