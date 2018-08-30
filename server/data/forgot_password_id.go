package data

// ForgotPasswordId: Generated when a user says they forgot their password
type ForgotPasswordId struct {
	Id     string  `gorm:"primary_key;unique;not null;size:190"`
	User   User    `gorm:"foreignKey:UserId"`
	UserId TUserID `gorm:"not null"`
	Used   bool    `gorm:"not null;default=false"`
}
