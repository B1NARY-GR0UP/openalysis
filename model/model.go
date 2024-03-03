package model

import "time"

// Model OPENALYSIS base model
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
