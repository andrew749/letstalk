package data

type AuthenticationData struct {
	UserId       int    `json:"user_id" gorm:"not null;primary_key;"`
	User         User   `gson:"foreignkey:UserId;association_foreignkey:UserId"`
	PasswordHash string `json:"password_hash" gorm:"not null;type:varchar(128);"`
}
