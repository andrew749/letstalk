package recommendations

import "letstalk/server/data"

type UserMatch struct {
	UserOneId data.TUserID
	UserTwoId data.TUserID
	Score     Score
}
