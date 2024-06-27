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

func CreateAnalyzedOrgContributions(ctx context.Context, db *gorm.DB, contributions map[string]model.AnalyzedOrgContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	var cs []model.AnalyzedOrgContribution
	for _, contribution := range contributions {
		cs = append(cs, contribution)
	}
	return db.WithContext(ctx).Create(&cs).Error
}

func CreateAnalyzedGroupContributions(ctx context.Context, db *gorm.DB, contributions map[string]model.AnalyzedGroupContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	var cs []model.AnalyzedGroupContribution
	for _, contribution := range contributions {
		cs = append(cs, contribution)
	}
	return db.WithContext(ctx).Create(&cs).Error
}
