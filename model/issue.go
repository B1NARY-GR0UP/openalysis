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
