package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model

	Name string

	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
	ContributorCount int
}
