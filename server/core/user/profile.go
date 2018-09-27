package user

import (
	"strconv"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

func ProfileEditController(c *ctx.Context) errs.Error {
	var request api.ProfileEditRequest
	err := c.GinContext.BindJSON(&request)
	if err != nil {
		return errs.NewRequestError(err.Error())
	}

	if request.Birthdate != nil {
		if requestErr := validateUserBirthday(*request.Birthdate); requestErr != nil {
			return requestErr
		}
	}

	if err := query.UpdateProfile(c.Db, c.Es, c.SessionData.UserId, request); err != nil {
		return err
	}
	return nil
}

func GetMyProfileController(c *ctx.Context) errs.Error {
	userModel, err := query.GetProfile(c.Db, c.SessionData.UserId, true)
	if err != nil {
		return nil
	}

	c.Result = *userModel
	return nil
}

func GetMatchProfileController(c *ctx.Context) errs.Error {
	matchUserIdStr := c.GinContext.Param("userId")
	tempMatchUserId, convErr := strconv.Atoi(matchUserIdStr)
	if convErr != nil {
		return errs.NewRequestError(convErr.Error())
	}

	matchUserId := data.TUserID(tempMatchUserId)

	userModel, err := query.GetMatchProfile(c.Db, c.SessionData.UserId, matchUserId)
	if err != nil {
		return err
	}

	c.Result = *userModel
	return nil
}

// GetPublicProfileController returns profile info when scanning a user's QR code.
func GetPublicProfileController(c *ctx.Context) errs.Error {
	userSecret := c.GinContext.Param("code")

	matchUser, dbErr := query.GetUserBySecret(c.Db, userSecret)
	if dbErr != nil {
		return errs.NewNotFoundError("no such user")
	}

	userModel, err := query.GetMatchProfile(c.Db, c.SessionData.UserId, matchUser.UserId)
	if err != nil {
		return err
	}

	c.Result = userModel
	return nil
}

func GetProfilePicUrl(ctx *ctx.Context) errs.Error {
	params := ctx.GinContext.Request.URL.Query()
	var userId int
	var val string

	if valTemp, ok := params["userId"]; ok {
		val = valTemp[0]
	} else {
		return errs.NewRequestError("Missing userId parameter")
	}

	if userIdTemp, err := strconv.Atoi(val); err == nil {
		userId = userIdTemp
	} else {
		return errs.NewRequestError("Malformed userId")
	}

	db := ctx.Db

	var user data.User
	if err := db.Select("profile_pic").Where("user_id = ?", userId).First(&user).Error; err != nil {
		return errs.NewInternalError(err.Error())
	}
	var profilePicResult api.ProfilePicResponse
	profilePicResult.ProfilePic = user.ProfilePic

	ctx.Result = profilePicResult
	return nil
}
