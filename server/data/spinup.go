package data

import (
	"letstalk/server/core/utility/uw_email"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
				// create required table
				err := tx.AutoMigrate(&UserDevice{}).Error

				// This method originally migrated notifications from sessions to a new table
				// we dont do that anymore and thus removed the logic for this (in the event of a new table spinup.)
				return err
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
		{
			ID: "Add other programs 2018 - 2023",
			Migrate: func(tx *gorm.DB) error {
				var programs = map[string]string{
					"ACCOUNTING_AND_FINANCIAL_MANAGEMENT":                        "Accounting and Financial Management",
					"ACTUARIAL_SCIENCE":                                          "Actuarial Science",
					"ANTHROPOLOGY":                                               "Anthropology",
					"APPLIED_MATHEMATICS":                                        "Applied Mathematics",
					"ARCHITECTURE":                                               "Architecture",
					"BIOCHEMISTRY":                                               "Biochemistry",
					"BIOLOGY":                                                    "Biology",
					"BIOMEDICAL_SCIENCES":                                        "Biomedical Sciences",
					"BIOSTATISTICS":                                              "Biostatistics",
					"BIOTECHNOLOGY/CHARTERED_PROFESSIONAL_ACCOUNTANCY":           "Biotechnology/Chartered Professional Accountancy",
					"BIOTECHNOLOGY/ECONOMICS":                                    "Biotechnology/Economics",
					"BUSINESS_ADMINISTRATION_AND_COMPUTER_SCIENCE_DOUBLE_DEGREE": "Business Administration and Computer Science Double Degree",
					"BUSINESS_ADMINISTRATION_AND_MATHEMATICS_DOUBLE_DEGREE":      "Business Administration and Mathematics Double Degree",
					"CHEMISTRY":                                 "Chemistry",
					"CLASSICAL_STUDIES":                         "Classical Studies",
					"COMBINATORICS_AND_OPTIMIZATION":            "Combinatorics and Optimization",
					"COMPUTATIONAL_MATHEMATICS":                 "Computational Mathematics",
					"COMPUTER_SCIENCE":                          "Computer Science",
					"COMPUTING_AND_FINANCIAL_MANAGEMENT":        "Computing and Financial Management",
					"DATA_SCIENCE":                              "Data Science",
					"EARTH_SCIENCES":                            "Earth Sciences",
					"ECONOMICS":                                 "Economics",
					"ENGLISH":                                   "English",
					"ENVIRONMENT_AND_BUSINESS":                  "Environment and Business",
					"ENVIRONMENT,_RESOURCES_AND_SUSTAINABILITY": "Environment, Resources and Sustainability",
					"ENVIRONMENTAL_SCIENCE":                     "Environmental Science",
					"FINE_ARTS":                                 "Fine Arts",
					"FRENCH":                                    "French",
					"GENDER_AND_SOCIAL_JUSTICE":                 "Gender and Social Justice",
					"GEOGRAPHY_AND_AVIATION":                    "Geography and Aviation",
					"GEOGRAPHY_AND_ENVIRONMENTAL_MANAGEMENT":    "Geography and Environmental Management",
					"GEOMATICS":                                 "Geomatics",
					"GERMAN":                                    "German",
					"GLOBAL_BUSINESS_AND_DIGITAL_ARTS":                   "Global Business and Digital Arts",
					"HEALTH_STUDIES":                                     "Health Studies",
					"HISTORY":                                            "History",
					"HONOURS_ARTS":                                       "Honours Arts",
					"HONOURS_ARTS_AND_BUSINESS":                          "Honours Arts and Business",
					"HONOURS_SCIENCE":                                    "Honours Science",
					"INFORMATION_TECHNOLOGY_MANAGEMENT":                  "Information Technology Management",
					"INTERNATIONAL_DEVELOPMENT":                          "International Development",
					"KINESIOLOGY":                                        "Kinesiology",
					"KNOWLEDGE_INTEGRATION":                              "Knowledge Integration",
					"LEGAL_STUDIES":                                      "Legal Studies",
					"LIBERAL_STUDIES":                                    "Liberal Studies",
					"LIFE_PHYSICS":                                       "Life Physics",
					"LIFE_SCIENCES":                                      "Life Sciences",
					"MATERIALS_AND_NANOSCIENCES":                         "Materials and Nanosciences",
					"MATHEMATICAL_ECONOMICS":                             "Mathematical Economics",
					"MATHEMATICAL_FINANCE":                               "Mathematical Finance",
					"MATHEMATICAL_OPTIMIZATION":                          "Mathematical Optimization",
					"MATHEMATICAL_PHYSICS":                               "Mathematical Physics",
					"MATHEMATICAL_STUDIES":                               "Mathematical Studies",
					"MATHEMATICS":                                        "Mathematics",
					"MATHEMATICS/BUSINESS_ADMINISTRATION":                "Mathematics/Business Administration",
					"MATHEMATICS/CHARTERED_PROFESSIONAL_ACCOUNTANCY":     "Mathematics/Chartered Professional Accountancy",
					"MATHEMATICS/FINANCIAL_ANALYSIS_AND_RISK_MANAGEMENT": "Mathematics/Financial Analysis and Risk Management",
					"MEDICINAL_CHEMISTRY":                                "Medicinal Chemistry",
					"MEDIEVAL_STUDIES":                                   "Medieval Studies",
					"MUSIC":                                              "Music",
					"PEACE_AND_CONFLICT_STUDIES":     "Peace and Conflict Studies",
					"PHILOSOPHY":                     "Philosophy",
					"PHYSICAL_SCIENCES":              "Physical Sciences",
					"PHYSICS":                        "Physics",
					"PHYSICS_AND_ASTRONOMY":          "Physics and Astronomy",
					"PLANNING":                       "Planning",
					"POLITICAL_SCIENCE":              "Political Science",
					"PSYCHOLOGY":                     "Psychology",
					"PUBLIC_HEALTH":                  "Public Health",
					"PURE_MATHEMATICS":               "Pure Mathematics",
					"RECREATION_AND_LEISURE_STUDIES": "Recreation and Leisure Studies",
					"RECREATION_AND_SPORT_BUSINESS":  "Recreation and Sport Business",
					"SCIENCE_AND_AVIATION":           "Science and Aviation",
					"SCIENCE_AND_BUSINESS":           "Science and Business",
					"SPANISH":                        "Spanish",
					"SPEECH_COMMUNICATION":           "Speech Communication",
					"STATISTICS":                     "Statistics",
					"THEATRE_AND_PERFORMANCE":        "Theatre and Performance",
					"THERAPEUTIC_RECREATION":         "Therapeutic Recreation",
					"TOURISM_DEVELOPMENT":            "Tourism Development",
					"OTHER":                          "Other",
				}

				for gradYear := uint(2018); gradYear <= uint(2023); gradYear++ {
					for programId, programName := range programs {
						sequenceName := "Other"
						sequenceId := "OTHER"
						cohort := &Cohort{
							ProgramId:    programId,
							ProgramName:  programName,
							GradYear:     gradYear,
							IsCoop:       false,
							SequenceName: &sequenceName,
							SequenceId:   &sequenceId,
						}
						err := db.Save(cohort).Error
						if err != nil {
							return err
						}
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add times to user table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add run id field to notification",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&Notification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add user_groups table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&UserGroup{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add campaigns",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(NotificationCampaign{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add surveys table",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate(&UserSurvey{})
				return tx.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Create jobmine",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(JobRecord{}).Error; err != nil {
					return err
				}
				return tx.AutoMigrate(TaskRecord{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add run_id to notifications table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(Notification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Backfill normalized @edu.uwaterloo.ca addresses",
			Migrate: func(tx *gorm.DB) error {
				rows, err := tx.Model(&VerifyEmailId{}).Rows()
				if err != nil {
					return err
				}
				defer rows.Close()

				for rows.Next() {
					verifyEmailId := VerifyEmailId{}
					if err := db.ScanRows(rows, &verifyEmailId); err != nil {
						return err
					}
					if !uw_email.Validate(verifyEmailId.Email) {
						return errors.Errorf("Could not normalize unexpected non-uw email '%s'", verifyEmailId.Email)
					}
					uwEmail := uw_email.FromString(verifyEmailId.Email)
					verifyEmailId.Email = uwEmail.ToStringNormalized()
					if err := db.Save(&verifyEmailId).Error; err != nil {
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
			ID: "Add other alum cohort 2005 - 2019",
			Migrate: func(tx *gorm.DB) error {
				sequenceName := "Other"
				sequenceId := "OTHER"
				for gradYear := uint(2005); gradYear <= uint(2019); gradYear++ {
					cohort := &Cohort{
						ProgramId:    "ALUM",
						ProgramName:  "Alum",
						GradYear:     gradYear,
						IsCoop:       false,
						SequenceName: &sequenceName,
						SequenceId:   &sequenceId,
					}
					err := db.Save(cohort).Error
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
			ID: "Create user_verify_link table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(UserVerifyLink{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add notification status updating job",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&ExpoPendingNotification{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&UserDevice{}).Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add meetup reminders table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(MeetupReminder{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add match round data models",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(MatchRound{}).Error; err != nil {
					return err
				}
				return tx.AutoMigrate(MatchRoundMatch{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "Add connection_match_round table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&ConnectionMatchRound{}).Error
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
