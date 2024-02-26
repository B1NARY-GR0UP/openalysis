package model

import "gorm.io/gorm"

type Repo struct {
	gorm.Model

	Name    string
	GroupID string // TODO: support multi groups

	Issue string
	PR    string

	Star  string
	Fork  string
	Clone string

	Contributor string
	// TODO
}
