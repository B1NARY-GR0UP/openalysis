package model

import "gorm.io/gorm"

type Cursor struct {
	gorm.Model
	EndCursor string
	// ISSUE | PR
	Type string
}
