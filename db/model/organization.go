package model

import "gorm.io/gorm"

// Organization one to many Repository
type Organization struct {
	gorm.Model

	Login  string
	NodeID string
}
