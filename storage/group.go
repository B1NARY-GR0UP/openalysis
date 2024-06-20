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

	"github.com/B1NARY-GR0UP/openalysis/model"
	"gorm.io/gorm"
)

func CreateGroup(ctx context.Context, db *gorm.DB, group *model.Group) error {
	return db.WithContext(ctx).Create(group).Error
}

func UpdateGroup(ctx context.Context, db *gorm.DB, group *model.Group) error {
	var currentGroup model.Group
	if err := db.WithContext(ctx).Where("name = ?", group.Name).First(&currentGroup).Error; err != nil {
		return err
	}
	currentGroup.IssueCount = group.IssueCount
	currentGroup.PullRequestCount = group.PullRequestCount
	currentGroup.StarCount = group.StarCount
	currentGroup.ForkCount = group.ForkCount
	currentGroup.ContributorCount = group.ContributorCount
	if err := db.WithContext(ctx).Save(&currentGroup).Error; err != nil {
		return err
	}
	return nil
}

func CreateGroupsOrganizations(ctx context.Context, db *gorm.DB, join *model.GroupsOrganizations) error {
	return db.WithContext(ctx).Create(join).Error
}

func CreateGroupsRepositories(ctx context.Context, db *gorm.DB, join *model.GroupsRepositories) error {
	return db.WithContext(ctx).Create(join).Error
}
