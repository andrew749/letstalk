package generic_notification_job

import (
	"database/sql"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/romana/rlog"
)

// GenericNotificationJob
// This job lets you send arbitrary notifications from jobmine (email and push).
// Currently the way to specify who to send notifications to is done via setting
// a metadata property USER_SELECTOR_QUERY that should return a row of the form:
//
// (user_id, otherProperties...)
//
// where other properties are arbitrary metadata that you want to pass to templating
// functions for Notifications and Email. The name of the column will be how you can
// access the property in the template.
//
// For example, for emails, this means if you return columns:
//  - user_id (remember you NEED to include this)
//  - first_name as first_name
//  - SUM(...) as num_mentees
// you will be able to access these properties in the template as user_id, first_name
// and num_mentees respectively.
//
// Similarly, the data parameter for job metadata will be available to every task in the
// same templating fashion.
//
// NOTE: A job will fail if there are duplicate keys in the data parameter and what gets
// selected from the database.

const GENERIC_NOTIFICATION_JOB jobmine.JobType = "GenericNotificationJob"

// How mysql sends back userId keys by default
const userIdDbKey = "user_id"

// rowToMap convert primary valued database columns to known types (from byte)
func rowToMap(columns []string, columnPointers []interface{}, columnValues []interface{}, rows *sql.Rows) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	err := rows.Scan(columnPointers...)
	if err != nil {
		return nil, err
	}
	for i, col := range columns {
		switch columnValues[i].(type) {
		case []byte:
			byteValue := columnValues[i].([]byte)
			res[col] = string(byteValue)
			rlog.Warnf("Processing byte array column %s => %s", col, res[col])
		case string:
			res[col] = columnValues[i].(string)
			rlog.Warnf("Processing string column %s => %s", col, res[col])
		case int:
			rlog.Warnf("Processing int column %s => %s", col, res[col])
		case int64:
			res[col] = columnValues[i].(int64)
			rlog.Warnf("Processing int64 column %s => %s", col, res[col])
		case float64:
			res[col] = columnValues[i].(float64)
			rlog.Warnf("Processing float column %s => %s", col, res[col])
		default:
			rlog.Warnf("Ignoring column %s => %v", col, columnValues[i])
		}
	}
	return res, nil
}

// getMetadataForQuery gets executes a query and packages as task metadata
func getMetadataForQuery(db *gorm.DB, query string) ([]TaskRecordMetadata, error) {
	// check for potential bad query
	if err := safetyCheck(query); err != nil {
		return nil, err
	}

	taskMetadata := make([]TaskRecordMetadata, 0)
	// execute sql query to get data for notifications
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	rlog.Debugf("Successfully executed user selector.")

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// For every column that we return per tuple, place this into a map to be passed
	// on to whatever templating function we're using downstream.
	columnPointers := make([]interface{}, len(columns))
	// backing store for data
	columnValues := make([]interface{}, len(columns))
	for i, _ := range columns {
		columnPointers[i] = &columnValues[i]
	}
	for rows.Next() {
		taskData, err := rowToMap(columns, columnPointers, columnValues, rows)
		if err != nil {
			return nil, err
		}

		userIdRaw, found := taskData[userIdDbKey]
		if !found {
			return nil, errors.New("Unable to get userId field. Terminating job.")
		}

		userIdInt, success := userIdRaw.(int64)
		if !success {
			return nil, errors.New("Unable to convert userId field to int64. Terminating job.")
		}

		taskMetadata = append(
			taskMetadata,
			TaskRecordMetadata{
				UserId: data.TUserID(uint(userIdInt)),
				Data:   taskData,
			},
		)
	}
	return taskMetadata, nil
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	jobRecordMetadata := parseJobRecordMetadata(jobRecord)

	rawTaskRecordMetadata, err := getMetadataForQuery(db, jobRecordMetadata.UserSelectorQuery)
	taskMetadata := make([]jobmine.Metadata, len(rawTaskRecordMetadata))
	for i, metadata := range rawTaskRecordMetadata {
		taskMetadata[i] = packageTaskRecordMetadata(metadata)
	}
	if err != nil {
		return nil, err
	}

	// don't write task records if dry run
	if jobRecordMetadata.DryRun {
		for _, metadata := range taskMetadata {
			rlog.Infof("Task Metadata: %+v", metadata)
		}
		return nil, nil
	}

	return taskMetadata, nil
}

// TODO: a safer solution is to create machinery that `switch`es onto a series of algorithms that
// can do specific user selections with type safety. This is obviously more flexible but it opens
// the opportunity for remote code execution. Hence why we have a safety method. A malicious user
// would need to gain access to our infra in this case but the safety check is still warranted.

// Sanity check to help prevent RCE
// rejects queries that create, insert, update, drop, delete, or alter a table. The idea is that we only allow selections.
func safetyCheck(query string) error {
	queryContains := func(phrase string) bool {
		return strings.Contains(strings.ToLower(query), phrase)
	}
	if queryContains("create") {
		return errors.New("Contains Create clause")
	}
	if queryContains("insert") {
		return errors.New("Contains Insert clause")
	}
	if queryContains("drop") {
		return errors.New("Contains Drop clause")
	}
	if queryContains("update") {
		return errors.New("Contains Update clause")
	}
	if queryContains("delete") {
		return errors.New("Contains Delete clause")
	}
	if queryContains("alter") {
		return errors.New("Contains Alter clause")
	}

	return nil
}

var GenericNotificationJob jobmine.JobSpec = jobmine.JobSpec{
	JobType:          GENERIC_NOTIFICATION_JOB,
	TaskSpec:         genericNotificationJobTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

// CreateGenericNotificationJob Creates a notification job record to get run at some point.
func CreateGenericNotificationJob(
	db *gorm.DB,
	runId string,
	dryRun bool,
	userSelectorQuery string,
	templatedTitle string,
	templatedMessage string,
	emailTemplateId *string,
	notificationTemplateId *string,
	additionalData map[string]interface{},
) error {
	metadata := packageJobRecordMetadata(userSelectorQuery, dryRun, templatedTitle, templatedMessage, emailTemplateId, notificationTemplateId, additionalData)

	safetyCheck(userSelectorQuery)

	if err := db.Create(&jobmine.JobRecord{
		JobType:  GENERIC_NOTIFICATION_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
