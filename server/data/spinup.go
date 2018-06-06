package data

import (
	"github.com/jinzhu/gorm"
)

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&AuthenticationData{})
	db.AutoMigrate(&Cohort{})
	PopulateCohort(db)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&UserVector{})
	db.AutoMigrate(&UserCohort{})
	db.AutoMigrate(&UserAdditionalData{})
	db.AutoMigrate(&NotificationToken{})
	db.AutoMigrate(&ExternalAuthData{})
	db.AutoMigrate(&FbAuthToken{})
	db.AutoMigrate(&Matching{})
	db.AutoMigrate(&RequestMatching{})
	db.AutoMigrate(&Credential{})
	db.AutoMigrate(&UserCredential{})
	db.AutoMigrate(&UserCredentialRequest{})
	db.AutoMigrate(&Subscriber{})
	db.AutoMigrate(&ForgotPasswordId{})
}

func CreateDB(db *gorm.DB) {
	migrateDB(db)
}
