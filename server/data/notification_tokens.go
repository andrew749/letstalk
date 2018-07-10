package data

type NotificationToken struct {
	SessionId string `gorm:"not null;primary_key"`
	Token     string `gorm:"not null;primary_key"`
	Service   string
}
