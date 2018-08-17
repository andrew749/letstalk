package data

import "time"

type Times struct {
	CreatedAt time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
}
