package recommendations

import (
	"time"

	"letstalk/server/data"
)

// 1 if user was created before Before, 0 otherwise.
// 0 for all users on blacklist.
type UserScoreOlder struct {
	Before           time.Time
	BlacklistUserIds map[data.TUserID]interface{}
}

func (s UserScoreOlder) RequiredObjects() []string {
	return []string{}
}

func (s UserScoreOlder) Calculate(user *data.User) (Score, error) {
	if _, ok := s.BlacklistUserIds[user.UserId]; ok {
		return 0.0, nil
	}
	if user.CreatedAt.Before(s.Before) {
		return 1.0, nil
	} else {
		return 0.0, nil
	}
}
