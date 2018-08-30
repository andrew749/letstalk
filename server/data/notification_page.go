package data

import "time"

// NotificationPage A notification page that is templated with user specific information.
type NotificationPage struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	User           User         `gorm:"foreign_key:UserId"`
	UserId         TUserID      `gorm:"primary_key;auto_increment:false"`
	Notification   Notification `gorm:"not null;foreign_key:NotificationId"`
	NotificationId uint         `gorm:"primary_key;auto_increment:false"`
	TemplateLink   string       `gorm:"not null;size:190"`
	Attributes     JSONBlob     `gorm:"not null;type:text"`
}
