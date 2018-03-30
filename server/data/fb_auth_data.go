package data

type FbAuthData struct {
	User     User    `gorm:"foreignkey:UserId"`
	UserId   *int    `json:"userId" gorm:"auto_increment;not null"`
	FbUserId *string `json:"fbUserId" gorm:"primary_key"`
}
