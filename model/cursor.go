package model

import (
	"gorm.io/gorm"
	"time"
)

type CursorType string

const (
	CursorTypeIssue CursorType = "ISSUE"
	CursorTypePR    CursorType = "PR"
)

type Cursor struct {
	gorm.Model
	RepoNodeID        string
	RepoNameWithOwner string
	// used by issue
	LastUpdate time.Time
	// used by pr
	EndCursor string
	// ISSUE | PR
	Type CursorType
}
