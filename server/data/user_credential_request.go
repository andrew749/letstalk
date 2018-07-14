package data

type UserCredentialRequest struct {
	User         User       `gorm:"foreignkey:UserId"`
	UserId       TUserID    `json:"userId" gorm:"not null,primary_key;auto_increment:false"`
	Credential   Credential `gorm:"foreignkey:CredentialId"`
	CredentialId uint       `gorm:"not null,primary_key;auto_increment:false"`
}
