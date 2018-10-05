package data

import (
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"gopkg.in/gormigrate.v1"
)

func isSQLite(db *gorm.DB) bool {
	return db.Dialect().GetName() == "sqlite3"
}

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
				// TODO: Do error checking here like below.
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
				return tx.AutoMigrate(&Notification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "3",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&NotificationPage{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "TRAITS_DATA_MODELS_V1_5",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&Organization{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&Role{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&UserPosition{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&SimpleTrait{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&UserSimpleTrait{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&UserLocation{}).Error
				if err != nil {
					return err
				}
				err = tx.AutoMigrate(&Cohort{}).Error
				if err != nil {
					return err
				}
				// NOTE: Need to make Cohort.SequenceId nullable, since we not longer enforce that it
				// exists.
				if !isSQLite(tx) {
					tx.Model(&Cohort{}).ModifyColumn("sequence_id", "varchar(190)")
				}
				return tx.AutoMigrate(&UserCohort{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Pending sent notifications",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&ExpoPendingNotification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "User Devices Creation From Session",
			Migrate: func(tx *gorm.DB) error {
				// for every session create a new user device entry
				if isSQLite(tx) {
					rlog.Warn("Skipping migration since sqlite")
					return nil
				}

				// create required table
				err := tx.AutoMigrate(&UserDevice{}).Error
				if err != nil {
					return err
				}

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
				return tx.AutoMigrate(SentMonthlyNotification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Verify email id",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&VerifyEmailId{}).Error
				if err != nil {
					return err
				}
				return tx.AutoMigrate(&User{}).Error // Added IsEmailVerified column.
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add deep linking field on notification",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&Notification{}).Error
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
		{
			ID: "Backfill soft/comp eng program/sequence names",
			Migrate: func(tx *gorm.DB) error {
				err := tx.Model(&Cohort{}).Where(&Cohort{ProgramId: "SOFTWARE_ENGINEERING"}).Update(&Cohort{
					ProgramName: "Software Engineering",
					IsCoop:      true,
				}).Error
				if err != nil {
					return err
				}

				err = tx.Model(&Cohort{}).Where(&Cohort{ProgramId: "COMPUTER_ENGINEERING"}).Update(&Cohort{
					ProgramName: "Computer Engineering",
					IsCoop:      true,
				}).Error
				if err != nil {
					return err
				}

				sequenceId := "4STREAM"
				sequenceName := "4 Stream"
				err = tx.Model(&Cohort{}).Where(&Cohort{SequenceId: &sequenceId}).Update(&Cohort{
					SequenceName: &sequenceName,
				}).Error
				if err != nil {
					return err
				}

				sequenceId = "8STREAM"
				sequenceName = "8 Stream"
				err = tx.Model(&Cohort{}).Where(&Cohort{SequenceId: &sequenceId}).Update(&Cohort{
					SequenceName: &sequenceName,
				}).Error
				if err != nil {
					return err
				}

				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add engineering cohorts 2018-2023",
			Migrate: func(tx *gorm.DB) error {
				var stream8Programs = map[string]string{
					"SOFTWARE_ENGINEERING":       "Software Engineering",
					"ELECTRICAL_ENGINEERING":     "Electrical Engineering",
					"COMPUTER_ENGINEERING":       "Computer Engineering",
					"CIVIL_ENGINEERING":          "Civil Engineering",
					"MANAGEMENT_ENGINEERING":     "Management Engineering",
					"NANOTECHNOLOGY_ENGINEERING": "Nanotechnology Engineering",
					"MECHANICAL_ENGINEERING":     "Mechanical Engineering",
					"MECHATRONICS_ENGINEERING":   "Mechatronics Engineering",
					"BIOMEDICAL_ENGINEERING":     "Biomedical Engineering",
					"CHEMICAL_ENGINEERING":       "Chemical Engineering",
				}
				var stream4Programs = map[string]string{
					"ARCHITECTURAL_ENGINEERING":  "Architectural Engineering",
					"ELECTRICAL_ENGINEERING":     "Electrical Engineering",
					"COMPUTER_ENGINEERING":       "Computer Engineering",
					"ENVIRONMENTAL_ENGINEERING":  "Environmental Engineering",
					"GEOLOGICAL_ENGINEERING":     "Geological Engineering",
					"SYSTEMS_DESIGN_ENGINEERING": "Systems Design Engineering",
					"MECHANICAL_ENGINEERING":     "Mechanical Engineering",
					"MECHATRONICS_ENGINEERING":   "Mechatronics Engineering",
					"CHEMICAL_ENGINEERING":       "Chemical Engineering",
				}

				for gradYear := uint(2018); gradYear <= uint(2023); gradYear++ {
					for programId, programName := range stream8Programs {
						sequenceName := "8 Stream"
						sequenceId := "8STREAM"
						cohort := &Cohort{
							ProgramId:    programId,
							ProgramName:  programName,
							GradYear:     gradYear,
							IsCoop:       true,
							SequenceName: &sequenceName,
							SequenceId:   &sequenceId,
						}
						err := db.Where(cohort).FirstOrCreate(cohort).Error
						if err != nil {
							return err
						}
					}

					for programId, programName := range stream4Programs {
						sequenceName := "4 Stream"
						sequenceId := "4STREAM"
						cohort := &Cohort{
							ProgramId:    programId,
							ProgramName:  programName,
							GradYear:     gradYear,
							IsCoop:       true,
							SequenceName: &sequenceName,
							SequenceId:   &sequenceId,
						}
						err := db.Where(cohort).FirstOrCreate(cohort).Error
						if err != nil {
							return err
						}
					}
				}

				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Make user birthdate optional",
			Migrate: func(tx *gorm.DB) error {
				// modify column to be nullable
				if !isSQLite(tx) {
					return db.Model(&User{}).ModifyColumn("birthdate", "varchar(100)").Error
				}
				return nil
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
