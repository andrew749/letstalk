package data

type UserCredential struct {
	User         User       `gorm:"foreignkey:UserId"`
	UserId       TUserID    `gorm:"not null,primary_key;auto_increment:false"`
	Credential   Credential `gorm:"foreignkey:CredentialId"`
	CredentialId uint       `gorm:"not null,primary_key;auto_increment:false"`
}
