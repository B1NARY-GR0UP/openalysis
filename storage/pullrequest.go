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

package storage

import (
	"context"
	"errors"

	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"gorm.io/gorm"
)

func CreatePullRequests(ctx context.Context, db *gorm.DB, prs []*model.PullRequest) error {
	if util.IsEmptySlice(prs) {
		return nil
	}
	return db.WithContext(ctx).Create(prs).Error
}

func UpdatePullRequest(ctx context.Context, db *gorm.DB, pr *model.PullRequest) error {
	var currentPR model.PullRequest
	if err := db.WithContext(ctx).Where("node_id = ?", pr.NodeID).First(&currentPR).Error; err != nil {
		return err
	}
	currentPR.State = pr.State
	currentPR.PRMergedAt = pr.PRMergedAt
	currentPR.PRClosedAt = pr.PRClosedAt
	if err := db.WithContext(ctx).Save(&currentPR).Error; err != nil {
		return err
	}
	return nil
}

func DeletePullRequests(ctx context.Context, db *gorm.DB, repoNodeID string) error {
	return db.WithContext(ctx).Where("repo_node_id = ?", repoNodeID).Delete(&model.PullRequest{}).Error
}

func QueryOPENPullRequests(ctx context.Context, db *gorm.DB, repoNodeID string) ([]model.PullRequest, error) {
	var prs []model.PullRequest
	err := db.WithContext(ctx).Where("state = ? AND repo_node_id = ?", "OPEN", repoNodeID).Find(&prs).Error
	return prs, err
}

func PullRequestAssigneesExist(ctx context.Context, db *gorm.DB, nodeID string) (bool, error) {
	var assignees model.PullRequestAssignees
	if err := db.WithContext(ctx).Where("pull_request_node_id = ?", nodeID).First(&assignees).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreatePullRequestAssignees(ctx context.Context, db *gorm.DB, assignees []*model.PullRequestAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return db.WithContext(ctx).Create(assignees).Error
}

func UpdatePullRequestAssignees(ctx context.Context, db *gorm.DB, prNodeID string, assignees []*model.PullRequestAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	var currentAssignees []*model.PullRequestAssignees
	if err := db.WithContext(ctx).Where("pull_request_node_id = ?", prNodeID).Find(&currentAssignees).Error; err != nil {
		return err
	}
	var s1 []model.PullRequestAssignees
	var s2 []model.PullRequestAssignees
	for _, assignee := range currentAssignees {
		s1 = append(s1, model.PullRequestAssignees{
			PullRequestNodeID:   assignee.PullRequestNodeID,
			PullRequestNumber:   assignee.PullRequestNumber,
			PullRequestURL:      assignee.PullRequestURL,
			PullRequestRepoName: assignee.PullRequestRepoName,
			AssigneeNodeID:      assignee.AssigneeNodeID,
			AssigneeLogin:       assignee.AssigneeLogin,
		})
	}
	for _, assignee := range assignees {
		s2 = append(s2, model.PullRequestAssignees{
			PullRequestNodeID:   assignee.PullRequestNodeID,
			PullRequestNumber:   assignee.PullRequestNumber,
			PullRequestURL:      assignee.PullRequestURL,
			PullRequestRepoName: assignee.PullRequestRepoName,
			AssigneeNodeID:      assignee.AssigneeNodeID,
			AssigneeLogin:       assignee.AssigneeLogin,
		})
	}
	more, less := util.CompareSlices(s1, s2)
	if !util.IsEmptySlice(more) {
		if err := db.WithContext(ctx).Create(more).Error; err != nil {
			return err
		}
	}
	for _, e := range less {
		if err := db.WithContext(ctx).Where("id = ?", e.ID).Delete(&model.PullRequestAssignees{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeletePullRequestAssigneesByPR(ctx context.Context, db *gorm.DB, prNodeID string) error {
	return db.WithContext(ctx).Where("pull_request_node_id = ?", prNodeID).Delete(&model.PullRequestAssignees{}).Error
}

func DeletePullRequestAssigneesByRepo(ctx context.Context, db *gorm.DB, nameWithOwner string) error {
	return db.WithContext(ctx).Where("pull_request_repo_name = ?", nameWithOwner).Delete(&model.PullRequestAssignees{}).Error
}
