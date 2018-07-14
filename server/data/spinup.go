package data

import (
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"gopkg.in/gormigrate.v1"
)

func migrateDB(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
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
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		rlog.Errorf("Could not migrate: %v", err)
	}

	rlog.Infof("Succesfully ran migration")
}

func CreateDB(db *gorm.DB) {
	migrateDB(db)
}
