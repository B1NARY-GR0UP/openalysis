package model

import (
	"time"
)

// PullRequest many to many Contributor
type PullRequest struct {
	Model

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
