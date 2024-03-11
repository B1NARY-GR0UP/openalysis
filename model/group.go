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

type GroupsRepositories struct {
	GroupName  string
	RepoNodeID string
}

type GroupsOrganizations struct {
	GroupName string
	OrgNodeID string
}
