package converters

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

func ApiUserPersonalInfoFromDataUser(user *data.User) api.UserPersonalInfo {
	return api.UserPersonalInfo{
		UserId:     user.UserId,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Gender:     user.Gender,
		Birthdate:  user.Birthdate,
		Secret:     user.Secret,
		ProfilePic: user.ProfilePic,
	}
}
