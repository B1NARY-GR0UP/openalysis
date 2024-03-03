package model

import "time"

// Model OPENALYSIS base model
// TODO: UpdatedAt and DeletedAt?? for groups update and delete
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
