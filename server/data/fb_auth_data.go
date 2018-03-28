package data

type FbAuthData struct {
	User     User    `gorm:"foreignkey:UserId"`
	UserId   *int    `json:"user_id" gorm:"auto_increment;not null"`
	FbUserId *string `json:"fb_user_id" gorm:"primary_key"`
}
