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

import "gorm.io/gorm"

// Contributor use github rest api
// Contributor one to many Issue (Author)
// Contributor one to many PullRequest (Author)
// Contributor many to many Issue (Assignees)
// Contributor many to many PullRequest (Assignees)
//
// NOTE: create on update
type Contributor struct {
	gorm.Model

	Login  string
	NodeID string

	Company   string
	Location  string
	AvatarURL string

	RepoOwner     string
	RepoName      string
	RepoNodeID    string
	Contributions int
}
