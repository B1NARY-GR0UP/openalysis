package db

import (
	"github.com/B1NARY-GR0UP/openalysis/model"
)

func CreateRepository(repo *model.Repository) (int64, error) {
	res := DB.Create(repo)
	return res.RowsAffected, res.Error
}
