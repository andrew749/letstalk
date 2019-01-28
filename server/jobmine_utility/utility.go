package jobmine_utility

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"letstalk/server/data"
	"letstalk/server/jobmine"
)

const TIME_LAYOUT = time.RFC3339

func TimeFromJobRecord(jobRecord jobmine.JobRecord, key string) (*time.Time, error) {
	if val, ok := jobRecord.Metadata[key]; ok {
		var (
			timeStr string
			ok      bool
		)
		if timeStr, ok = val.(string); !ok {
			return nil, errors.New(fmt.Sprintf("%s must be a time string", key))
		}
		time, err := time.Parse(TIME_LAYOUT, timeStr)
		if err != nil {
			return nil, err
		}
		return &time, nil
	}
	return nil, nil
}

func UIntFromJobRecord(jobRecord jobmine.JobRecord, key string) (*uint, error) {
	if valIntf, exists := jobRecord.Metadata[key]; exists {
		if valFloat, isFloat := valIntf.(float64); isFloat {
			val := uint(valFloat)
			return &val, nil
		} else {
			return nil, errors.New(fmt.Sprintf("%s must be a number, got %v", key, valIntf))
		}
	}
	return nil, nil
}

func FormatTime(tme time.Time) string {
	return tme.Format(TIME_LAYOUT)
}

func UserIdFromTaskRecord(taskRecord jobmine.TaskRecord, key string) (*data.TUserID, error) {
	if val, ok := taskRecord.Metadata[key]; ok {
		var (
			userIdFloat float64
			ok          bool
		)
		if userIdFloat, ok = val.(float64); !ok {
			return nil, errors.New(fmt.Sprintf("%s must be a userId (number), but got %v", key, val))
		}
		userId := data.TUserID(uint(userIdFloat))
		return &userId, nil
	}
	return nil, nil
}
