package data

type NotificationToken struct {
	SessionId string `gorm:"not null;primary_key;size:100"`
	Token     string `gorm:"not null;primary_key;size:100"`
	Service   string
}
