package db

import (
	"github.com/B1NARY-GR0UP/openalysis/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// TODO: replace dsn with config
	var err error
	DB, err = gorm.Open(mysql.Open("root:114514@tcp(localhost:3306)/openalysis?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// TODO: use mount
	DB.AutoMigrate(&model.Repository{})
}

func CreateRepository(repo *model.Repository) (int64, error) {
	res := DB.Create(repo)
	return res.RowsAffected, res.Error
}
