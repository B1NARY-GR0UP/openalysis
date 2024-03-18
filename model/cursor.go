package model

import (
	"gorm.io/gorm"
	"time"
)

type Cursor struct {
	gorm.Model
	RepoNodeID        string
	RepoNameWithOwner string
	// used by issue
	LastUpdate time.Time
	// used by pr
	EndCursor string
}
