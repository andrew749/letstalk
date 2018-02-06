package login

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/db"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"time"

	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
)

/**
 * The basic structure of a request to create a new user.
 * The following data must be POSTed to create a new user.
 * Users cannot have the same email address.

  {
    'first_name': string,
    'last_name': string,
    'email': string,
		'phone_number' string(optional),
		'gender': string,
		'birthday': date,
    'password': string,
  }
*/

/**
 * Parse the json request and bind the parameters to a struct.
 */
func getUserDataFromRequest(c *ctx.Context) (*api.User, errs.Error) {
	var inputUser api.User
	err := c.GinContext.BindJSON(&inputUser)
	if err != nil {
		return nil, errs.NewClientError("%s", err)
	}
	rlog.Debugf("post user: %s", inputUser)
	return &inputUser, nil
}

func SignupUser(c *ctx.Context) errs.Error {
	// get the data that the user submitted in the post request
	user, err := getUserDataFromRequest(c)

	if err != nil {
		return err
	}

	// Check that no user exists with this email.
	existingUser, dberr := data.UserObjs.Select().Where(data.UserObjs.FilterEmail("=", user.Email)).List(c.Db)

	if dberr != nil {
		return errs.NewDbError(err)
	}

	if len(existingUser) != 0 {
		return errs.NewClientError("a user already exists with email: %s", user.Email)
	}

	return writeUser(user, c)
}

/**
 * Create a new user given a particular request and insert in the db.
 */
func writeUser(user *api.User, c *ctx.Context) errs.Error {
	// Create user data structures in the orm.

	userModel := data.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    utility.GenderIdByName(user.Gender),
		Birthdate: time.Unix(user.Birthday, 0),
	}

	var err error

	if userModel.UserId, err = db.NumId(c); err != nil {
		return errs.NewDbError(err)
	}

	hashedPassword, err := utility.HashPassword(*user.Password)

	if err != nil {
		return errs.NewInternalError("Unable to hash password")
	}

	authData := data.AuthenticationData{
		userModel.UserId,
		hashedPassword,
	}

	// Insert data structures within a transaction.
	dbErr := gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
		if _, err := userModel.Insert(tx); err != nil {
			return err
		}
		if _, err := authData.Insert(tx); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	user.Password = nil
	c.Result = user
	return nil
}
