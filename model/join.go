package model

type GroupsRepositories struct {
	GroupID    uint
	RepoNodeID string
}

type GroupsOrganizations struct {
	GroupID   uint
	OrgNodeID string
}

// IssueAssignees an issue can have multi assignees
type IssueAssignees struct {
	IssueNodeID    string
	AssigneeNodeID string
	AssigneeLogin  string
}

// PullRequestAssignees a pr can have multi reviewers
type PullRequestAssignees struct {
	PullRequestNodeID string
	AssigneeNodeID    string
	AssigneeLogin     string
}
