package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Company  string
	Location string
	// TODO
}
