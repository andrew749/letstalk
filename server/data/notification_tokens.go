package data

type NotificationToken struct {
	Token   string `json:"token" gorm:"not null;primary_key"`
	Service string `json:"service"`
}
