package data

type ExternalAuthData struct {
	User        User    `gorm:"foreignkey:UserId"`
	UserId      int     `gorm:"primary_key;not null"`
	FbUserId    *string `gorm:"null"`
	PhoneNumber *string `gorm:"null"`
}
