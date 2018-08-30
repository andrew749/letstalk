package data

import "github.com/jinzhu/gorm"

type SentMonthlyNotification struct {
	gorm.Model
	Matching   Matching `gorm:"foreign_key:MatchingId;"`
	MatchingId uint     `gorm:"unique_index:monthly_unique_index"`
	RunId      string   `gorm:"unique_index:monthly_unique_index"` // id of the run that generated this notification
}

func SentOutMonthlyNotification(db *gorm.DB, matchingId uint, runId string) error {
	monthlyNotification := SentMonthlyNotification{
		MatchingId: matchingId,
		RunId:      runId,
	}
	if err := db.FirstOrCreate(&monthlyNotification).Error; err != nil {
		return err
	}
	return nil
}
