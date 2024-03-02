package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model

	Name string
}
