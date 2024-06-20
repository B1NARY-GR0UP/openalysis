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

func CreateIssues(ctx context.Context, db *gorm.DB, issues []*model.Issue) error {
	if util.IsEmptySlice(issues) {
		return nil
	}
	return db.WithContext(ctx).Create(issues).Error
}

func IssueExist(ctx context.Context, db *gorm.DB, nodeID string) (bool, error) {
	var issue model.Issue
	if err := db.WithContext(ctx).Where("node_id = ?", nodeID).First(&issue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateIssue(ctx context.Context, db *gorm.DB, issue *model.Issue) error {
	var currentIssue model.Issue
	if err := db.WithContext(ctx).Where("node_id = ?", issue.NodeID).First(&currentIssue).Error; err != nil {
		return err
	}
	currentIssue.State = issue.State
	currentIssue.IssueClosedAt = issue.IssueClosedAt
	if err := db.WithContext(ctx).Save(&currentIssue).Error; err != nil {
		return err
	}
	return nil
}

func DeleteIssues(ctx context.Context, db *gorm.DB, repoNodeID string) error {
	return db.WithContext(ctx).Where("repo_node_id = ?", repoNodeID).Delete(&model.Issue{}).Error
}

func CreateIssueAssignees(ctx context.Context, db *gorm.DB, assignees []*model.IssueAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	return db.WithContext(ctx).Create(assignees).Error
}

func IssueAssigneesExist(ctx context.Context, db *gorm.DB, nodeID string) (bool, error) {
	var assignees model.IssueAssignees
	if err := db.WithContext(ctx).Where("issue_node_id = ?", nodeID).First(&assignees).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateIssueAssignees(ctx context.Context, db *gorm.DB, issueNodeID string, assignees []*model.IssueAssignees) error {
	if util.IsEmptySlice(assignees) {
		return nil
	}
	var currentAssignees []*model.IssueAssignees
	if err := db.WithContext(ctx).Where("issue_node_id = ?", issueNodeID).Find(&currentAssignees).Error; err != nil {
		return err
	}
	var s1 []model.IssueAssignees
	var s2 []model.IssueAssignees
	for _, assignee := range currentAssignees {
		s1 = append(s1, model.IssueAssignees{
			IssueNodeID:    assignee.IssueNodeID,
			IssueNumber:    assignee.IssueNumber,
			IssueURL:       assignee.IssueURL,
			IssueRepoName:  assignee.IssueRepoName,
			AssigneeNodeID: assignee.AssigneeNodeID,
			AssigneeLogin:  assignee.AssigneeLogin,
		})
	}
	for _, assignee := range assignees {
		s2 = append(s2, model.IssueAssignees{
			IssueNodeID:    assignee.IssueNodeID,
			IssueNumber:    assignee.IssueNumber,
			IssueURL:       assignee.IssueURL,
			IssueRepoName:  assignee.IssueRepoName,
			AssigneeNodeID: assignee.AssigneeNodeID,
			AssigneeLogin:  assignee.AssigneeLogin,
		})
	}
	more, less := util.CompareSlices(s1, s2)
	if !util.IsEmptySlice(more) {
		if err := db.WithContext(ctx).Create(more).Error; err != nil {
			return err
		}
	}
	for _, e := range less {
		if err := db.WithContext(ctx).Where("id = ?", e.ID).Delete(&model.IssueAssignees{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteIssueAssigneesByIssue(ctx context.Context, db *gorm.DB, issueNodeID string) error {
	return db.WithContext(ctx).Where("issue_node_id = ?", issueNodeID).Delete(&model.IssueAssignees{}).Error
}

func DeleteIssueAssigneesByRepo(ctx context.Context, db *gorm.DB, nameWithOwner string) error {
	return db.WithContext(ctx).Where("issue_repo_name = ?", nameWithOwner).Delete(&model.IssueAssignees{}).Error
}
