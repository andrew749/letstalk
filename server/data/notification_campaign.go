package data

import (
	"bytes"

	"github.com/jinzhu/gorm"
)

type NotificationCampaign struct {
	gorm.Model
	RunId         string  `gorm:"primary_key;size:190"`
	FailedUserIds *string `gorm:"text"`
}

func (c *NotificationCampaign) SetFailedUserIds(db *gorm.DB, failedIds []TUserID) error {
	var buffer bytes.Buffer
	for i, userId := range failedIds {
		buffer.Write([]byte(string(userId)))
		// keep adding commas
		if i != len(failedIds)-1 {
			buffer.Write([]byte(","))
		}
	}
	failedIdsString := buffer.String()
	c.FailedUserIds = &failedIdsString
	return db.Save(c).Error
}

// ExistsCampaign Check if a campaign exists
func ExistsCampaign(db *gorm.DB, RunId string) (bool, error) {
	var campaign NotificationCampaign
	if err := db.Where(&NotificationCampaign{RunId: RunId}).First(&campaign).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateCampaign Creates a new campaign
func CreateCampaign(db *gorm.DB, RunId string) (*NotificationCampaign, error) {
	campaign := NotificationCampaign{RunId: RunId}
	if err := db.Create(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}
