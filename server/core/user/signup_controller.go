package user

import (
	"bytes"
	"encoding/base64"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/s3_assets"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"letstalk/server/email"

	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
func getUserDataFromRequest(c *ctx.Context) (*api.SignupRequest, error) {
	var inputUser api.SignupRequest
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
		return errs.NewRequestError(err.Error())
	}

	err = c.Db.Model(&data.User{}).Where("email = ?", user.Email).First(&data.User{}).Error
	if err == nil {
		return errs.NewRequestError("a user already exists with email: %s", user.Email)
	} else if err != nil && !gorm.IsRecordNotFoundError(err) { // Some other db error
		return errs.NewDbError(err)
	}

	if user.Birthdate != nil {
		if requestErr := validateUserBirthday(*user.Birthdate); requestErr != nil {
			return requestErr
		}
	}

	err = writeUser(user, c)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	err = email.SendNewAccountEmail(
		mail.NewEmail(user.FirstName, user.Email),
		user.FirstName,
	)

	// don't fail if we can't send an email
	if err != nil {
		raven.CaptureError(err, nil)
		rlog.Error(err)
	}

	return nil
}

// Birthday must be in YYYY-MM-DD format.
func validateUserBirthday(birthday string) errs.Error {
	birthdate, err := time.Parse(utility.BirthdateFormat, birthday)
	if err != nil {
		return errs.NewRequestError("Bad user birthday format")
	}
	if utility.Today().AddDate(-13, 0, 0).Before(birthdate) {
		return errs.NewRequestError("Must be at least 13 years old")
	}
	return nil
}

/**
 * Create a new user given a particular request and insert in the db.
 */
func writeUser(userData *api.SignupRequest, c *ctx.Context) error {
	// Create user data structures in the orm.

	// Insert data structures within a transaction.
	tx := c.Db.Begin()
	var userModel *data.User
	var err error
	if userModel, err = CreateUserWithAuth(
		tx,
		userData.Email,
		userData.FirstName,
		userData.LastName,
		userData.Gender,
		userData.Birthdate,
		data.USER_ROLE_DEFAULT,
		userData.Password,
	); err != nil {
		tx.Rollback()
		return err
	}

	externalAuthRecord := data.ExternalAuthData{
		UserId:      userModel.UserId,
		PhoneNumber: &userData.PhoneNumber,
	}

	if err := tx.Create(&externalAuthRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	// upload the profile pic
	if userData.ProfilePic != nil {
		var photoData []byte
		if photoData, err = base64.StdEncoding.DecodeString(*userData.ProfilePic); err != nil {
			return err
		}
		reader := bytes.NewReader(photoData)
		var location *string
		if location, err = s3_assets.UploadProfilePic(userModel.UserId, reader); err != nil {
			tx.Rollback()
			return err
		}

		rlog.Debug("Successfully uploaded profile pic. Updating profile pic")
		if err = tx.Model(&userModel).Update(&data.User{ProfilePic: location}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	c.Result = struct{ UserId data.TUserID }{userModel.UserId}
	return nil
}
