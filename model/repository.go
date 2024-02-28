package model

import "gorm.io/gorm"

type Repository struct {
	gorm.Model

	Owner string
	Name  string

	GroupID string // TODO: support multi groups

	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int

	// TODO
}
