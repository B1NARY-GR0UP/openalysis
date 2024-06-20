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

func CreateOrganization(ctx context.Context, db *gorm.DB, org *model.Organization) error {
	return db.WithContext(ctx).Create(org).Error
}

func UpdateOrganization(ctx context.Context, db *gorm.DB, org *model.Organization) error {
	var currentOrg model.Organization
	if err := db.WithContext(ctx).Where("node_id = ?", org.NodeID).First(&currentOrg).Error; err != nil {
		return err
	}
	currentOrg.IssueCount = org.IssueCount
	currentOrg.PullRequestCount = org.PullRequestCount
	currentOrg.StarCount = org.StarCount
	currentOrg.ForkCount = org.ForkCount
	currentOrg.ContributorCount = org.ContributorCount
	if err := db.WithContext(ctx).Save(&currentOrg).Error; err != nil {
		return err
	}
	return nil
}
