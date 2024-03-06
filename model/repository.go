package model

import "gorm.io/gorm"

// Repository one to many Contributor
type Repository struct {
	gorm.Model

	Owner  string // might be a user or org
	Name   string
	NodeID string

	OwnerNodeID string

	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
	ContributorCount int
}
