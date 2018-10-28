package notifications

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/pkg/errors"
	"github.com/romana/rlog"

	"github.com/jinzhu/gorm"
)

// SendNotificationCampaign send a campaign to  a group of users.
func SendNotificationCampaign(db *gorm.DB, req api.NotificationCampaignSendRequest) error {
	campaignExists, err := data.ExistsCampaign(db, req.RunId)
	if err != nil {
		return err
	}

	// send notification to each user
	if campaignExists {
		rlog.Debugf("Campaign with id %s already exists", req.RunId)
		return errors.New(fmt.Sprintf("Campaign already exists %+v", err))
	}

	var campaign *data.NotificationCampaign
	if campaign, err = data.CreateCampaign(db, req.RunId); err != nil {
		rlog.Debugf("Unable to create campaign %s because %+v", req.RunId, err)
		return errors.New(fmt.Sprintf("Unable to create new campaign %+v", err))
	}

	// the campaign is fresh, lets send stuff out
	users, err := query.GetUsersByGroupId(db, data.TGroupID(req.GroupId))
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to get users by groupId %+v", req.GroupId))
	}

	var ce *errs.CompositeError
	var failedUserIds = make([]data.TUserID, 0)
	for _, user := range users {
		rlog.Debugf("[Campaign %s]: Trying to send notification to user %d", req.RunId, user.UserId)
		// TODO(acod): add templating to the title, message and metadata
		var (
			userId       = user.UserId
			title        = req.Title
			message      = req.Message
			thumbnail    = req.Thumbnail
			templatePath = req.TemplatePath
			metadata     = req.TemplatedMetadata
			runId        = req.RunId
		)

		// inject some additional metadata to make it easier to personalize notifications
		metadata["user"] = user

		if err := CreateAdHocNotification(db, userId, title, message, thumbnail, templatePath, metadata, &runId); err != nil {
			rlog.Errorf("Unable to create notification for user %d because %+v", userId, err)
			failedUserIds = append(failedUserIds, userId)
			ce = errs.AppendNullableError(ce, err)
		}
	}

	// if there were any failed user ids to send
	if len(failedUserIds) > 0 {
		if err := campaign.SetFailedUserIds(db, failedUserIds); err != nil {
			rlog.Error("Failed to update which user Ids were unable to send")
			ce = errs.AppendNullableError(ce, err)
		}
	}
	if ce != nil {
		return ce
	}
	return nil
}
