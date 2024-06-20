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
	"gorm.io/gorm"
)

func CreateCursor(ctx context.Context, db *gorm.DB, cursor *model.Cursor) error {
	return db.WithContext(ctx).Create(cursor).Error
}

func QueryCursor(ctx context.Context, db *gorm.DB, repo string) (*model.Cursor, error) {
	cursor := &model.Cursor{}
	err := db.WithContext(ctx).Where("repo_name_with_owner = ?", repo).First(cursor).Error
	// for organization's new repository case
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cursor, nil
	}
	return cursor, err
}

func UpdateOrCreateCursor(ctx context.Context, db *gorm.DB, cursor *model.Cursor) error {
	var currentCursor model.Cursor
	if err := db.WithContext(ctx).Where("repo_node_id = ?", cursor.RepoNodeID).First(&currentCursor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&model.Cursor{
				RepoNodeID:        cursor.RepoNodeID,
				RepoNameWithOwner: cursor.RepoNameWithOwner,
				LastUpdate:        cursor.LastUpdate,
				EndCursor:         cursor.EndCursor,
			}).Error; err != nil {
				return err
			}
			return nil
		}
		return err
	}
	currentCursor.RepoNameWithOwner = cursor.RepoNameWithOwner
	currentCursor.LastUpdate = cursor.LastUpdate
	currentCursor.EndCursor = cursor.EndCursor
	if err := db.WithContext(ctx).Save(&currentCursor).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCursor(ctx context.Context, db *gorm.DB, repoNodeID string) error {
	return db.WithContext(ctx).Where("repo_node_id = ?", repoNodeID).Delete(&model.Cursor{}).Error
}
