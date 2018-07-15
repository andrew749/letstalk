package contact_info

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"strconv"

	"github.com/jinzhu/gorm"
)

/**
 * Get request required the requested userId data as a url parameter
 * ?userId=<>
 */
func GetContactInfoController(c *ctx.Context) errs.Error {
	params := c.GinContext.Request.URL.Query()
	var userId data.TUserID
	var val string
	if valTemp, ok := params["userId"]; ok {
		val = valTemp[0]
	} else {
		return errs.NewRequestError("Missing userId parameter")
	}

	if userIdTemp, err := strconv.Atoi(val); err == nil {
		userId = data.TUserID(userIdTemp)
	} else {
		return errs.NewRequestError("Malformed userId")
	}

	if res, err := isAllowedToAccessContactInfo(
		c.Db,
		c.SessionData.UserId,
		userId,
	); err == nil && res == true {
		user, err := query.GetUserById(c.Db, userId)
		if err != nil {
			if _, ok := err.(*errs.NotFoundError); ok == true {
				return errs.NewNotFoundError("Unable to get user: %s", err)
			}

			// otherwise internal error
			return errs.NewInternalError(err.Error())
		}
		c.Result = api.ContactInfo{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
	} else if err != nil {
		return errs.NewInternalError(err.Error())
	} else {
		return errs.NewForbiddenError("Not allowed to access this user's contact info")
	}
	return nil
}

/**
 * Determine if a user is allowed to access specific information
 */
func isAllowedToAccessContactInfo(db *gorm.DB, requestorId data.TUserID, requestedId data.TUserID) (bool, error) {
	var count int
	// find who this person is a mentor for
	if err := db.Table("mentors").
		Where("mentor_id = ? AND user_user_id = ?", requestorId, requestedId).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	if err := db.Table("mentees").
		Where("mentee_id = ? AND user_user_id = ?", requestorId, requestedId).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}
