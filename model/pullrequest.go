package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	PullRequestStateTypeOpen   = "OPEN"
	PullRequestStateTypeMerged = "MERGED"
	PullRequestStateTypeClosed = "CLOSED"
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

	PullRequestNodeID   string
	PullRequestNumber   int
	PullRequestURL      string
	PullRequestRepoName string

	AssigneeNodeID string
	AssigneeLogin  string
}
