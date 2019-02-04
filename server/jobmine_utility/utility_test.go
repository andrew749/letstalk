package jobmine_utility

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"letstalk/server/jobmine"
)

func TestTimeFromJobRecord(t *testing.T) {
	now := time.Now().UTC()
	sec, err := time.ParseDuration("1s")
	assert.NoError(t, err)
	now = now.Truncate(sec)
	metadata := make(map[string]interface{})
	metadata["time"] = FormatTime(now)
	jobRecord := jobmine.JobRecord{Metadata: metadata}

	newNow, err := TimeFromJobRecord(jobRecord, "time")
	assert.NoError(t, err)
	assert.Equal(t, now, *newNow)
}

func TestTimeFromJobRecordMissing(t *testing.T) {
	metadata := make(map[string]interface{})
	jobRecord := jobmine.JobRecord{Metadata: metadata}
	newNow, err := TimeFromJobRecord(jobRecord, "time")
	assert.NoError(t, err)
	assert.Nil(t, newNow)
}

func TestTimeFromJobRecordInvalidFormat(t *testing.T) {
	metadata := make(map[string]interface{})
	metadata["time"] = "not a time"
	jobRecord := jobmine.JobRecord{Metadata: metadata}

	newNow, err := TimeFromJobRecord(jobRecord, "time")
	assert.Error(t, err)
	assert.Nil(t, newNow)
}
