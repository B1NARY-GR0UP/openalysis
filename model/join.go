package model

type GroupsRepositories struct {
	GroupName  string
	RepoNodeID string
}

type GroupsOrganizations struct {
	GroupName string
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
