package generic_notification_job

import (
	"letstalk/server/data"
	"letstalk/server/jobmine"
)

// Task Record Metadata
const (
	USER_ID_METADATA_KEY = "userId"
	TASK_DATA_KEY        = "data"
)

type TaskRecordMetadata struct {
	UserId data.TUserID
	Data   map[string]interface{}
}

func packageTaskRecordMetadata(taskRecordMetadata TaskRecordMetadata) map[string]interface{} {
	return map[string]interface{}{
		USER_ID_METADATA_KEY: taskRecordMetadata.UserId,
		TASK_DATA_KEY:        taskRecordMetadata.Data,
	}
}

func parseTaskRecordMetadata(taskRecord jobmine.TaskRecord) TaskRecordMetadata {
	return TaskRecordMetadata{
		UserId: data.TUserID(uint(taskRecord.Metadata[USER_ID_METADATA_KEY].(float64))),
		Data:   taskRecord.Metadata[TASK_DATA_KEY].(map[string]interface{}),
	}
}
