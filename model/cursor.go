package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	CursorTypeIssue = "ISSUE"
	CursorTypePR    = "PR"
)

type Cursor struct {
	gorm.Model
	RepoNodeID string
	LastUpdate time.Time
	// ISSUE | PR
	Type string
}
