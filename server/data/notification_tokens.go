package data

type NotificationToken struct {
	SessionId string `gorm:"not null;primary_key;size:190"`
	Token     string `gorm:"not null;primary_key;size:190"`
	Service   string
}
