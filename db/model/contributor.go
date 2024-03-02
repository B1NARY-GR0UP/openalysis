package model

import "gorm.io/gorm"

// Contributor use github rest api
// Contributor one to many Issue (Author)
// Contributor one to many PullRequest (Author)
// Contributor many to many Issue (Assignees)
// Contributor many to many PullRequest (Assignees)
type Contributor struct {
	gorm.Model

	Login  string
	NodeID string

	Company   string
	Location  string
	AvatarURL string

	RepoNodeID    string
	Contributions int
}
