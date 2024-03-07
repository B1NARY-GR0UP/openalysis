package model

import "gorm.io/gorm"

type Cursor struct {
	gorm.Model
	IssueEndCursor string
	PREndCursor    string
}
