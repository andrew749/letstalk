package data

type UserAdditionalData struct {
	User                 User    `gorm:"foreignkey:UserId"`
	UserId               TUserID `json:"userId" gorm:"not null;primary_key;auto_increment:false"`
	MentorshipPreference *int
	Bio                  *string `gorm:"type:text"`
	Hometown             *string
}
