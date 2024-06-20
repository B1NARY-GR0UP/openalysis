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

func CreateRepository(ctx context.Context, db *gorm.DB, repo *model.Repository) error {
	return db.WithContext(ctx).Create(repo).Error
}

func QueryRepositoryNodeID(ctx context.Context, db *gorm.DB, owner, name string) (string, error) {
	var repo model.Repository
	err := db.WithContext(ctx).Where(model.Repository{
		Owner: owner,
		Name:  name,
	}).First(&repo).Error
	return repo.NodeID, err
}

func DeleteRepository(ctx context.Context, db *gorm.DB, nodeID string) error {
	return db.WithContext(ctx).Where("node_id = ?", nodeID).Delete(&model.Repository{}).Error
}

func QueryReposByOrg(ctx context.Context, db *gorm.DB, orgNodeID string) ([]string, error) {
	var repos []model.Repository
	if err := db.WithContext(ctx).Where("owner_node_id = ?", orgNodeID).Group("node_id").Find(&repos).Error; err != nil {
		return nil, err
	}
	var reposNameWithOwner []string
	for _, repo := range repos {
		reposNameWithOwner = append(reposNameWithOwner, util.MergeNameWithOwner(repo.Owner, repo.Name))
	}
	return reposNameWithOwner, nil
}
