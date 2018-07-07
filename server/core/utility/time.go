package utility

import "time"

const (
	BirthdateFormat = "2006-01-02"
)

func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
