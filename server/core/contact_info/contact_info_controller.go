package contact_info

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"strconv"

	"github.com/mijia/modelq/gmq"
)

type ContactInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

/**
 * Get request required the requested userId data as a url parameter
 */
func GetContactInfoController(c *ctx.Context) errs.Error {
	params := c.GinContext.Request.URL.Query()
	var userId int
	var val string
	if valTemp, ok := params["userId"]; ok {
		val = valTemp[0]
	} else {
		return errs.NewClientError("Missing userId parameter")
	}

	if userIdTemp, err := strconv.Atoi(val); err == nil {
		userId = userIdTemp
	} else {
		return errs.NewClientError("Malformed userId")
	}

	if res, err := isAllowedToAccessContactInfo(
		c.Db,
		c.SessionData.UserId,
		userId,
	); err == nil && res == true {
		user, err := api.GetUserWithId(c.Db, userId)
		if err != nil {
			return errs.NewClientError("Unable to get user: %s", err)
		}
		c.Result = ContactInfo{user.FirstName, user.LastName, user.Email}
	} else if err != nil {
		return errs.NewInternalError(err.Error())
	} else {
		return errs.NewClientError("Not allowed to access this user's contact info")
	}
	return nil
}

/**
 * Determine if a user is allowed to access specific information
 */
func isAllowedToAccessContactInfo(db *gmq.Db, requestorId int, requestedId int) (bool, error) {
	// TODO write tests for this.
	// check if the user making the request is the persons mentor or mentee
	stmt, err := db.Prepare(`
		SELECT COUNT(*)
		FROM matchings
		WHERE
			(matchings.mentor=?  AND matchings.mentee=?) OR
			(matchings.mentor=?  AND matchings.mentee=?)
	`)

	if err != nil {
		return false, err
	}

	res, err := stmt.Query(
		requestorId,
		requestedId,
		requestedId,
		requestorId,
	)

	if err != nil {
		return false, err
	}

	// see if there are any matchings between this user and the requested user
	var count int
	res.Next()
	err = res.Scan(&count)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
