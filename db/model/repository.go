package model

// Repository TODO: 每次更新都需要插入新的 item，但是 NodeID 是一样的，ID 是自增的，结合 CreatedAt 来绘制时间序列图
// Repository one to many Contributor
type Repository struct {
	Model

	Owner  string // might be a user or org
	Name   string
	NodeID string

	OwnerNodeID string

	IssueCount       int
	PullRequestCount int
	StarCount        int
	ForkCount        int
}
