package jobmine

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const TIME_LAYOUT = time.RFC3339

func TimeFromJobRecord(jobRecord JobRecord, key string) (*time.Time, error) {
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

func FormatTime(tme time.Time) string {
	return tme.Format(TIME_LAYOUT)
}
