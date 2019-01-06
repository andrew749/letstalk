package notification_status_checker

import (
	"letstalk/server/data"
	"letstalk/server/jobmine"
	notification_api "letstalk/server/notifications"
	"letstalk/server/queue/queues/notification_queue"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const NOTIFICATION_STATUS_CHECKER_JOB jobmine.JobType = "NotificationStatusChecker"

const (
	NOTIFICATION_ID_KEY              = "notificationId"
	EXPO_PENDING_NOTIFICATION_ID_KEY = "expoPendingNotificationId"
)

type taskRecordMetadata struct {
	notificationId            uint
	expoPendingNotificationId uint
}

func packageTaskRecordMetadata(recordMetaData taskRecordMetadata) map[string]interface{} {
	return map[string]interface{}{
		NOTIFICATION_ID_KEY:              recordMetaData.notificationId,
		EXPO_PENDING_NOTIFICATION_ID_KEY: recordMetaData.expoPendingNotificationId,
	}
}

func parseTaskRecordData(data jobmine.Metadata) taskRecordMetadata {
	taskRecordMetadata := taskRecordMetadata{
		notificationId:            uint(data[NOTIFICATION_ID_KEY].(float64)),
		expoPendingNotificationId: uint(data[EXPO_PENDING_NOTIFICATION_ID_KEY].(float64)),
	}
	return taskRecordMetadata
}

func processNotification(db *gorm.DB, e *data.ExpoPendingNotification) error {
	// check the status of this notification from the server
	serverStatus, err := notification_api.GetNotificationStatus([]string{*e.Receipt})
	if err != nil {
		return err
	}

	status := serverStatus.Data[*e.Receipt].Status
	// if the status is ok, mark the notification as such
	if status == notification_api.OK_STATUS {
		return e.MarkNotificationDelivered(db)
	}

	// if there is an error update the state accordingly and try to remediate
	failureType := notification_api.ExpoNotificationFailureType(serverStatus.Data[*e.Receipt].Details.Error)

	// check the status of the notification on expo's side and perform remediation
	switch failureType {
	case notification_api.ErrorDeviceNotRegistered:
		// remove the device token from the database
		rlog.Errorf("Device registration is not valid anymore.")
		err := data.RemoveUserDevice(db, e.DeviceId)
		if err != nil {
			return err
		}
		break
	case notification_api.ErrorMessageTooBig:
		rlog.Errorf("Message is too big to send.")
		break
	case notification_api.ErrorMessageRateExceeded:
		rlog.Errorf("Message rate exceeded")
		notification, err := data.GetNotification(db, e.NotificationId)
		if err != nil {
			return err
		}
		if err := notification_queue.PushNotificationToQueue(*notification); err != nil {
			return err
		}
		break
	default:
		// wtf is happening?
		rlog.Errorf("Unknown error: %+v", serverStatus.Data[*e.Receipt].Details.Error)
	}
	return nil
}

var NotificationStatusChecker jobmine.JobSpec = jobmine.JobSpec{
	JobType: NOTIFICATION_STATUS_CHECKER_JOB,
	TaskSpec: jobmine.TaskSpec{
		Execute: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord) (interface{}, error) {
			rlog.Infof("Got data from taskRecord %s", taskRecord.Metadata["key"])
			taskRecordMetadata := parseTaskRecordData(taskRecord.Metadata)

			// get the actual notification object
			pendingNotification, err := data.GetPendingNotification(db, taskRecordMetadata.expoPendingNotificationId)
			if err != nil {
				return nil, err
			}

			return nil, processNotification(db, pendingNotification)
		},
		OnError: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
			taskRecordMetadata := parseTaskRecordData(taskRecord.Metadata)
			rlog.Errorf("Failed to update status of notification: %#v\n%+v", taskRecordMetadata, err)
		},
		OnSuccess: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, res interface{}) {
			taskRecordMetadata := parseTaskRecordData(taskRecord.Metadata)
			rlog.Infof("Successfully updated status of notification: %#v", taskRecordMetadata)
		},
	},
	GetTasksToCreate: func(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
		res := make([]jobmine.Metadata, 0)
		var notifications []data.ExpoPendingNotification

		// Find all pending notifications without success or error state and check
		if err := db.Where("notification_status", nil).Find(&notifications).Error; err != nil {
			return nil, err
		}

		// parse messages and package into jobmine format
		for _, notification := range notifications {
			metadata := taskRecordMetadata{
				notificationId:            notification.NotificationId,
				expoPendingNotificationId: notification.ID,
			}
			res = append(res, jobmine.Metadata(packageTaskRecordMetadata(metadata)))
		}

		return res, nil
	},
}