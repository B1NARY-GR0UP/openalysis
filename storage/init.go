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
	"log"
	"os"
	"time"

	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() (err error) {
	dsn := util.AssembleDSN(
		config.GlobalConfig.DataSource.MySQL.Host,
		config.GlobalConfig.DataSource.MySQL.Port,
		config.GlobalConfig.DataSource.MySQL.User,
		config.GlobalConfig.DataSource.MySQL.Password,
		config.GlobalConfig.DataSource.MySQL.Database,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true, // turn on ignore record not found error option
			Colorful:                  true,
		}),
	})
	if err != nil {
		return
	}

	err = DB.AutoMigrate(
		&model.Cursor{},
		&model.Contributor{},
		&model.Group{},
		&model.Issue{},
		&model.Organization{},
		&model.PullRequest{},
		&model.Repository{},
		&model.GroupsOrganizations{},
		&model.GroupsRepositories{},
		&model.IssueAssignees{},
		&model.PullRequestAssignees{},
	)
	if err != nil {
		return
	}
	return
}
