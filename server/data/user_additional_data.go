package data

type UserAdditionalData struct {
	User                 User `gorm:"foreignkey:UserId"`
	UserId               int  `json:"userId" gorm:"not null;primary_key;auto_increment:false"`
	MentorshipPreference *int
	Bio                  *string
	Hometown             *string
}
