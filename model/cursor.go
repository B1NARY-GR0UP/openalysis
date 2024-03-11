package model

import "gorm.io/gorm"

const (
	CursorTypeIssue = "ISSUE"
	CursorTypePR    = "PR"
)

type Cursor struct {
	gorm.Model
	RepoNodeID string
	EndCursor  string
	// ISSUE | PR
	Type string
}
