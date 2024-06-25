// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"time"

	"gorm.io/gorm"
)

// PullRequest many to many Contributor
type PullRequest struct {
	gorm.Model

	NodeID string

	Author       string
	AuthorNodeID string
	RepoNodeID   string
	RepoOwner    string
	RepoName     string
	Number       int

	// CLOSED | MERGED | OPEN
	State string

	PRCreatedAt time.Time
	PRMergedAt  *time.Time
	PRClosedAt  *time.Time
}

// PullRequestAssignees a pr can have multi reviewers
// a user can be assigned to multi pull requests
type PullRequestAssignees struct {
	gorm.Model

	PullRequestNodeID string
	PullRequestNumber int
	PullRequestURL    string
	// repo name with owner
	PullRequestRepoName string

	AssigneeNodeID string
	AssigneeLogin  string
}
