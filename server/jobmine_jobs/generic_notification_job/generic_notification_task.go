package generic_notification_job

import (
	"bytes"
	"html/template"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/email"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func mergeMaps(d1 map[string]interface{}, d2 map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(d1)+len(d2))
	for key, value := range d1 {
		result[key] = value
	}
	for key, value := range d2 {
		if _, exists := result[key]; exists {
			return nil, errs.NewBaseError("Unable to merge maps. Duplicate key %s", key)
		}
		result[key] = value
	}
	return result, nil
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	taskRecordMetadata := parseTaskRecordMetadata(taskRecord)
	jobRecordMetadata := parseJobRecordMetadata(jobRecord)
	var user data.User

	if err := db.Where("user_id = ? ", taskRecordMetadata.UserId).First(&user).Error; err != nil {
		return nil, err
	}
	var (
		title          string
		message        string
		templateBuffer bytes.Buffer
	)

	mergedTemplateData, err := mergeMaps(taskRecordMetadata.Data, jobRecordMetadata.AdditionalData)
	if err != nil {
		return nil, err
	}

	titleTemplate, err := template.New("titleTemplate").Parse(jobRecordMetadata.Title)
	if err != nil {
		return nil, err
	}
	err = titleTemplate.Execute(&templateBuffer, mergedTemplateData)
	if err != nil {
		return nil, nil
	}

	title = templateBuffer.String()
	templateBuffer.Reset()

	messageTemplate, err := template.New("messageTemplate").Parse(jobRecordMetadata.Message)
	err = messageTemplate.Execute(&templateBuffer, mergedTemplateData)
	if err != nil {
		return nil, nil
	}
	message = templateBuffer.String()
	templateBuffer.Reset()
	// send push notification
	if jobRecordMetadata.NotificationTemplate != nil {
		err := notifications.CreateAdHocNotificationNoTransaction(
			db,
			taskRecordMetadata.UserId,
			title,
			message,
			nil,
			*jobRecordMetadata.NotificationTemplate,
			taskRecordMetadata.Data,
			&jobRecord.RunId,
		)
		if err != nil {
			return nil, err
		}

	}

	// send email
	if jobRecordMetadata.EmailTemplate != nil {
		to := mail.NewEmail(user.FirstName, user.Email)
		if err := email.SendBasicTemplatedEmailFromMap(
			to,
			*jobRecordMetadata.EmailTemplate,
			mergedTemplateData,
		); err != nil {
			return nil, err
		}
	}

	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	taskRecordMetadata := parseTaskRecordMetadata(taskRecord)
	rlog.Infof("Unable to send message to user with id=%d: %+v", taskRecordMetadata.UserId, err)
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	taskRecordMetadata := parseTaskRecordMetadata(taskRecord)
	rlog.Infof("Successfully sent message to user with id=%d", taskRecordMetadata.UserId)
}

var genericNotificationJobTaskSpec = jobmine.TaskSpec{
	Execute:   execute,
	OnError:   onError,
	OnSuccess: onSuccess,
}
