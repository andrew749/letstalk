package data

import "time"

// Include this in data models that require the different time fields. Note that when entities
// containing this struct are deleted, they remain in the DB, but their DeletedAt field is non-null.
type Times struct {
	CreatedAt time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
}
