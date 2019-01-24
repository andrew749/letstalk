package verify_link

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/utility"
)

const LINK_TYPE_WHITELIST_WINTER_2019 data.VerifyLinkType = "WHITELIST_WINTER_2019"

// Creates a UserVerifyLink and returns the id for the link
func CreateLink(
	db *gorm.DB,
	userId data.TUserID,
	linkType data.VerifyLinkType,
	expiresAt *time.Time,
) (*data.TVerifyLinkID, errs.Error) {
	link := &data.UserVerifyLink{
		Id:        data.TVerifyLinkID(uuid.New().String()),
		UserId:    userId,
		Clicked:   false,
		Type:      linkType,
		ExpiresAt: expiresAt,
	}
	err := db.Save(link).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	return &link.Id, nil
}

func ClickLink(db *gorm.DB, linkId data.TVerifyLinkID) errs.Error {
	var link data.UserVerifyLink

	err := utility.WrapGormDBError(
		db.Where(&data.UserVerifyLink{Id: linkId}).First(&link).Error,
		"Link not found",
	)
	if err != nil {
		return err
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return errs.NewRequestError("Link is expired")
	}

	link.Clicked = true
	if err := db.Save(&link).Error; err != nil {
		return errs.NewDbError(err)
	}

	return nil
}
