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
	"github.com/B1NARY-GR0UP/openalysis/util"
	"gorm.io/gorm"
)

func CreateContributors(ctx context.Context, db *gorm.DB, cs []*model.Contributor) error {
	if util.IsEmptySlice(cs) {
		return nil
	}
	return db.WithContext(ctx).Create(cs).Error
}

// UpdateContributorCompanyAndLocation
// TODO: use batch update
func UpdateContributorCompanyAndLocation(ctx context.Context, db *gorm.DB, update func(string) string) error {
	var contributors []model.Contributor
	if err := db.WithContext(ctx).Find(&contributors).Error; err != nil {
		return err
	}
	for _, contributor := range contributors {
		contributor.Company = update(contributor.Company)
		contributor.Location = update(contributor.Location)
		if err := db.WithContext(ctx).Save(&contributor).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateContributorCompanyAndLocationByLogin
// TODO: use batch update
func UpdateContributorCompanyAndLocationByLogin(ctx context.Context, db *gorm.DB, login, company, location string) error {
	var currentContributors []model.Contributor
	if err := db.WithContext(ctx).Where("login = ?", login).Find(&currentContributors).Error; err != nil {
		return err
	}
	// if empty string, then do not update
	// TODO: optimize, according to SRP
	for _, contributor := range currentContributors {
		if company != "" {
			contributor.Company = company
		}
		if location != "" {
			contributor.Location = location
		}
		if err := db.WithContext(ctx).Save(&contributor).Error; err != nil {
			return err
		}
	}
	return nil
}

// QueryContributorCountByOrg
//
// SELECT COUNT(DISTINCT c.node_id) AS contributor_count
// FROM contributors c
// INNER JOIN repositories r ON c.repo_node_id = r.node_id
// WHERE r.owner_node_id = 'orgNodeID';
func QueryContributorCountByOrg(ctx context.Context, db *gorm.DB, orgNodeID string) (int, error) {
	var contributorCount int
	if err := db.WithContext(ctx).
		Table("contributors").
		Select("COUNT(DISTINCT contributors.node_id) AS contributor_count").
		Joins("INNER JOIN repositories ON contributors.repo_node_id = repositories.node_id").
		Where("repositories.owner_node_id = ?", orgNodeID).
		Scan(&contributorCount).Error; err != nil {
		return 0, err
	}
	return contributorCount, nil
}

// QueryContributorCountByGroup
//
// SELECT COUNT(DISTINCT c.node_id) AS contributor_count
// FROM contributors c
// INNER JOIN (
//
//	SELECT DISTINCT gr.repo_node_id
//	FROM groups_repositories gr
//	INNER JOIN repositories r ON gr.repo_node_id = r.node_id
//	WHERE gr.group_name = 'groupName'
//	UNION
//	SELECT DISTINCT r.node_id
//	FROM repositories r
//	INNER JOIN groups_organizations go ON r.owner_node_id = go.org_node_id
//	WHERE go.group_name = 'groupName'
//
// ) AS repos ON c.repo_node_id = repos.repo_node_id;
func QueryContributorCountByGroup(ctx context.Context, db *gorm.DB, groupName string) (int, error) {
	var count int64

	var repos1 []string
	sq1 := db.WithContext(ctx).
		Table("groups_repositories").
		Select("groups_repositories.repo_node_id").
		Joins("INNER JOIN repositories ON groups_repositories.repo_node_id = repositories.node_id").
		Where("groups_repositories.group_name = ?", groupName)
	if err := sq1.Find(&repos1).Error; err != nil {
		return 0, err
	}

	var repos2 []string
	sq2 := db.WithContext(ctx).
		Table("repositories").
		Select("repositories.node_id").
		Joins("INNER JOIN groups_organizations ON repositories.owner_node_id = groups_organizations.org_node_id").
		Where("groups_organizations.group_name = ?", groupName)
	if err := sq2.Find(&repos2).Error; err != nil {
		return 0, err
	}

	repoNodeIDs := append(repos1, repos2...)

	if err := db.WithContext(ctx).
		Table("contributors").
		Select("contributors.node_id").
		Where("contributors.repo_node_id IN ?", repoNodeIDs).
		Distinct().
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func UpdateOrCreateContributors(ctx context.Context, db *gorm.DB, cs []*model.Contributor) error {
	for _, contributor := range cs {
		if err := db.WithContext(ctx).Where(model.Contributor{
			NodeID:     contributor.NodeID,
			RepoNodeID: contributor.RepoNodeID,
		}).Assign(contributor).FirstOrCreate(contributor).Error; err != nil {
			return err
		}
	}
	return nil
}
