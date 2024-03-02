package model

type GroupsRepositories struct {
	GroupID    uint
	RepoNodeID string
}

type GroupsOrganizations struct {
	GroupID   uint
	OrgNodeID string
}

type ContributorsIssues struct {
	ContribNodeID string
	IssueNodeID   string
}

type ContributorsPullRequests struct {
	ContribNodeID string
	PRNodeID      string
}
