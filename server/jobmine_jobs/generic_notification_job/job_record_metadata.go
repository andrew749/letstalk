package generic_notification_job

import (
	"letstalk/server/jobmine"

	"github.com/romana/rlog"
)

// Job Record metadata
const (
	DRY_RUN                  = "dryRun"
	USER_SELECTOR_QUERY      = "userSelectorQuery"
	TITLE_METADATA_KEY       = "title"
	MESSAGE_METADATA_KEY     = "message"
	EMAIL_TEMPLATE_ID        = "emailTemplate"
	NOTIFICATION_TEMPLATE_ID = "notificationTemplate"
	ADDTIONAL_DATA           = "data"
)

type JobRecordMetadata struct {
	DryRun               bool
	UserSelectorQuery    string
	Title                string
	Message              string
	EmailTemplate        *string
	NotificationTemplate *string
	AdditionalData       map[string]interface{}
}

func packageJobRecordMetadata(
	userSelectorQuery string,
	dryRun bool,
	title string,
	message string,
	emailTemplate *string,
	notificationTemplate *string,
	additionalData map[string]interface{},
) map[string]interface{} {
	return map[string]interface{}{
		DRY_RUN:                  dryRun,
		USER_SELECTOR_QUERY:      userSelectorQuery,
		TITLE_METADATA_KEY:       title,
		MESSAGE_METADATA_KEY:     message,
		EMAIL_TEMPLATE_ID:        emailTemplate,
		NOTIFICATION_TEMPLATE_ID: notificationTemplate,
		ADDTIONAL_DATA:           additionalData,
	}
}

func parseJobRecordMetadata(
	jobRecord jobmine.JobRecord,
) JobRecordMetadata {
	var (
		emailTemplate        *string
		notificationTemplate *string
	)
	rlog.Debugf("%+v", jobRecord.Metadata)

	if rawEmailTemplate, found := jobRecord.Metadata[EMAIL_TEMPLATE_ID]; found {
		rlog.Debugf("Found email template: %s", rawEmailTemplate)
		tempEmailTemplate := rawEmailTemplate.(string)
		emailTemplate = &tempEmailTemplate
	}

	if rawNotificationTemplate, found := jobRecord.Metadata[NOTIFICATION_TEMPLATE_ID]; found {
		rlog.Debugf("Found notification template: %s", rawNotificationTemplate)
		tempNotificationTemplate := rawNotificationTemplate.(string)
		notificationTemplate = &tempNotificationTemplate
	}

	return JobRecordMetadata{
		DryRun:               jobRecord.Metadata[DRY_RUN].(bool),
		UserSelectorQuery:    jobRecord.Metadata[USER_SELECTOR_QUERY].(string),
		Title:                jobRecord.Metadata[TITLE_METADATA_KEY].(string),
		Message:              jobRecord.Metadata[MESSAGE_METADATA_KEY].(string),
		EmailTemplate:        emailTemplate,
		NotificationTemplate: notificationTemplate,
		AdditionalData:       jobRecord.Metadata[ADDTIONAL_DATA].(map[string]interface{}),
	}
}
