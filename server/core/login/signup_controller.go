package login

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

/**
 * The basic structure of a request to create a new user.
 * The following data must be POSTed to create a new user.
 * Users cannot have the same email address.

 Request:
  {
    'first_name': string,
    'last_name': string,
    'email': string,
		'phone_number' string(optional),
		'gender': string,
		'birthday': date,
    'password': string,
  }


	Response
	{
		user_id: string,
  }
*/

/**
 * Parse the json request and bind the parameters to a struct.
 */
func getUserDataFromRequest(c *ctx.Context) (*api.User, error) {
	var inputUser api.User
	err := c.GinContext.BindJSON(&inputUser)
	if err != nil {
		return nil, err
	}
	rlog.Debugf("post user: %s", inputUser)
	return &inputUser, nil
}

func SignupUser(c *ctx.Context) errs.Error {
	// get the data that the user submitted in the post request
	user, err := getUserDataFromRequest(c)

	if err != nil {
		return errs.NewClientError(err.Error())
	}

	err = c.Db.Model(&data.User{}).Where("email = ?", user.Email).First(&data.User{}).Error
	if err == nil {
		return errs.NewClientError("a user already exists with email: %s", user.Email)
	} else if err != nil && !gorm.IsRecordNotFoundError(err) { // Some other db error
		return errs.NewDbError(err)
	}

	err = writeUser(user, c)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

/**
 * Create a new user given a particular request and insert in the db.
 */
func writeUser(user *api.User, c *ctx.Context) error {
	// Create user data structures in the orm.

	bday := time.Unix(user.Birthday, 0)

	userModel := data.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    utility.GenderIdByName(user.Gender),
		Birthdate: &bday,
	}

	var err error

	hashedPassword, err := utility.HashPassword(*user.Password)

	if err != nil {
		return errs.NewInternalError("Unable to hash password")
	}

	authData := data.AuthenticationData{
		UserId:       userModel.UserId,
		PasswordHash: hashedPassword,
	}

	// Insert data structures within a transaction.
	tx := c.Db.Begin()
	if err := tx.Create(&userModel).Error; err != nil {
		tx.Rollback()
		return err
	}
	authData.UserId = userModel.UserId
	if err := tx.Create(&authData).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	c.Result = struct{ UserId int }{userModel.UserId}
	return nil
}
