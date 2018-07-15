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
				tx.AutoMigrate(&AuthenticationData{})
				tx.AutoMigrate(&Cohort{})
				PopulateCohort(tx)
				tx.AutoMigrate(&User{})
				tx.AutoMigrate(&Session{})
				tx.AutoMigrate(&UserVector{})
				tx.AutoMigrate(&UserCohort{})
				tx.AutoMigrate(&UserAdditionalData{})
				tx.AutoMigrate(&NotificationToken{})
				tx.AutoMigrate(&ExternalAuthData{})
				tx.AutoMigrate(&FbAuthToken{})
				tx.AutoMigrate(&Matching{})
				tx.AutoMigrate(&RequestMatching{})
				tx.AutoMigrate(&Credential{})
				tx.AutoMigrate(&UserCredential{})
				tx.AutoMigrate(&UserCredentialRequest{})
				tx.AutoMigrate(&Subscriber{})
				tx.AutoMigrate(&ForgotPasswordId{})
				tx.AutoMigrate(&MeetingConfirmation{})
				return tx.Error
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
