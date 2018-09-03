package data

import (
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"gopkg.in/gormigrate.v1"
)

func migrateDB(db *gorm.DB) {
	options := gormigrate.Options{
		TableName:      "migrations",
		IDColumnName:   "id",
		IDColumnSize:   190,
		UseTransaction: false,
	}
	m := gormigrate.New(db, &options, []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&AuthenticationData{})
				tx.AutoMigrate(&Cohort{})
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
		{
			ID: "2",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&Notification{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "3",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&NotificationPage{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "TRAITS_DATA_MODELS_V1_5",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&Organization{})
				tx.AutoMigrate(&Role{})
				tx.AutoMigrate(&UserPosition{})
				tx.AutoMigrate(&SimpleTrait{})
				tx.AutoMigrate(&UserSimpleTrait{})
				tx.AutoMigrate(&UserLocation{})
				tx.AutoMigrate(&Cohort{})
				// NOTE: Need to make Cohort.SequenceId nullable, since we not longer enforce that it
				// exists.
				tx.Model(&Cohort{}).ModifyColumn("sequence_id", "varchar(190)")
				tx.AutoMigrate(&UserCohort{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Pending sent notifications",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&ExpoPendingNotification{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "User Devices Creation From Session",
			Migrate: func(tx *gorm.DB) error {
				// for every session create a new user device entry

				// create required table
				tx.AutoMigrate(&UserDevice{})

				// row to scan results into
				type Row struct {
					token string
					uid   uint
				}

				rows, err := tx.Table("notification_tokens").
					Select("notification_tokens.token, sessions.user_id").
					Joins("inner join sessions on sessions.session_id=notification_tokens.session_id").
					Rows()

				if err != nil {
					return err
				}

				for rows.Next() {
					res := Row{}

					err := rows.Scan(&res.token, &res.uid)
					if err != nil {
						return err
					}

					// insert row into devices table
					err = AddExpoDeviceTokenforUser(tx, TUserID(res.uid), res.token)

					if err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add book-keeping for monthly notification",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(SentMonthlyNotification{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Verify email id",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&VerifyEmailId{})
				tx.AutoMigrate(&User{}) // Added IsEmailVerified column.
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "user_connections_v1",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&Connection{})
				tx.AutoMigrate(&ConnectionIntent{})
				tx.AutoMigrate(&Mentorship{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		rlog.Errorf("Could not migrate: %v", err)
		panic(err)
	}

	rlog.Infof("Succesfully ran migration")
}

func CreateDB(db *gorm.DB) {
	migrateDB(db)
}
