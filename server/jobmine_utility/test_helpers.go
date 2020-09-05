package jobmine_utility

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/jobmine"
)

func RunAndTestRunners(t *testing.T, db *gorm.DB, runId string, specStore jobmine.JobSpecStore) {
	now := time.Now()
	var queryJob jobmine.JobRecord
	err := db.Where("run_id = ?", runId).Find(&queryJob).Error
	assert.NoError(t, err)
	assert.NotZero(t, queryJob.ID)
	assert.Equal(t, jobmine.STATUS_CREATED, queryJob.Status)

	// set start time a little earlier
	queryJob.StartTime = now.AddDate(0, -1, 0) // last month
	err = db.Save(&queryJob).Error
	assert.NoError(t, err)

	// schedule tasks
	err = jobmine.JobRunner(specStore, db, nil, nil)
	assert.NoError(t, err)

	// run tasks
	err = jobmine.TaskRunner(specStore, db, nil, nil)
	assert.NoError(t, err)

	// update the state of all tasks
	err = jobmine.JobStateWatcher(db)
	assert.NoError(t, err)

	// check the state of all jobs
	err = db.Where("run_id = ?", runId).First(&queryJob).Error
	assert.NoError(t, err)
	assert.Equal(t, queryJob.Status, jobmine.STATUS_SUCCESS)

	// check the state of tasks
	var queryTasks []jobmine.TaskRecord
	err = db.Where(&jobmine.TaskRecord{JobId: queryJob.ID}).Find(&queryTasks).Error
	assert.NoError(t, err)
	for _, queryTask := range queryTasks {
		assert.Equal(t, jobmine.STATUS_SUCCESS, queryTask.Status)
	}
}
