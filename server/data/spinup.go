package data

import (
	"github.com/jinzhu/gorm"
)

func migrateDB(db *gorm.DB) {
	db = db.Set("gorm:table_options", "CHARSET=utf8mb4") // Create all new tables with utf8mb4 encoding.
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
	db.AutoMigrate(&MeetingConfirmation{})
}

func CreateDB(db *gorm.DB) {
	migrateDB(db)
}
