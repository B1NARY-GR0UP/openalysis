package model

import "gorm.io/gorm"

// Organization one to many Repository
type Organization struct {
	gorm.Model

	Login  string
	NodeID string
	// TODO: name (public profile name)

	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
	ContributorCount int
}
