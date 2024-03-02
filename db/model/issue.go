package model

import (
	"time"
)

// Issue many to many Contributor
type Issue struct {
	Model

	NodeID string

	// TODO: author could be bot, user...
	Author       string
	AuthorNodeID string
	RepoNodeID   string
	Number       int

	// OPEN | CLOSED
	State string

	IssueCreatedAt time.Time
	IssueClosedAt  time.Time
}
