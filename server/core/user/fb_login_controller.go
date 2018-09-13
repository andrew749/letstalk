package user

/**
 * Controller to handle user logging into
 */

import (
	"errors"
	"fmt"
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"letstalk/server/email"

	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
	fb "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func FBController(c *ctx.Context) errs.Error {
	var loginRequest api.FBLoginRequestData
	var externalAuthRecord data.ExternalAuthData
	var userId data.TUserID

	err := c.GinContext.BindJSON(&loginRequest)

	if err != nil {
		return errs.NewRequestError("%s", err)
	}

	authToken := loginRequest.Token
	expiry := time.Unix(loginRequest.Expiry, 0)

	user, err := getFBUser(authToken)

	if err != nil {
		return errs.NewRequestError("%s", err)
	}

	tx := c.Db.Begin()

	// check if the user already has facebook
	if tx.Where("fb_user_id = ?", user.Id).First(&externalAuthRecord).RecordNotFound() {

		appUser := data.User{
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Gender:          user.Gender,
			Birthdate:       user.Birthdate,
			Role:            data.USER_ROLE_DEFAULT,
			IsEmailVerified: false,
			ProfilePic:      &user.ProfilePic,
		}

		// Generate UUID for FB user.
		secret, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback()
			return errs.NewInternalError("%v", err)
		}
		appUser.Secret = secret.String()

		if err := tx.Where(&appUser).FirstOrCreate(&appUser).Error; err != nil {
			tx.Rollback()
			rlog.Error("Unable to insert new user")
			return errs.NewRequestError("Unable to create user")
		}

		userId = appUser.UserId
		rlog.Debug("Created new user with id: ", userId)

		externalAuthRecord.FbUserId = &user.Id
		externalAuthRecord.UserId = userId
		// insert the user's fb auth data
		if err := tx.Create(&externalAuthRecord).Error; err != nil {
			tx.Rollback()
			rlog.Error(err)
			return errs.NewRequestError("Unable to create user")
		}

		// the user Id for this application

		fbAuthToken := data.FbAuthToken{
			UserId:    userId,
			AuthToken: authToken,
			Expiry:    expiry,
		}
		if err := tx.Create(&fbAuthToken).Error; err != nil {
			tx.Rollback()
			rlog.Error(err)
			return errs.NewDbError(err)
		}

		if err := tx.Commit().Error; err != nil {
			rlog.Error(err)
			return errs.NewDbError(err)
		}
		//send email
		if err := email.SendNewAccountEmail(
			mail.NewEmail(appUser.FirstName, appUser.Email),
			appUser.FirstName,
		); err != nil {
			raven.CaptureError(err, nil)
			rlog.Error(err)
		}
	} else {
		rlog.Infof("Already found facebook user with fbid %#v", externalAuthRecord)
		userId = externalAuthRecord.UserId
	}

	// create new session for user id
	session, err := (*c.SessionManager).CreateNewSessionForUserId(userId, &loginRequest.NotificationToken)

	if err != nil {
		rlog.Error("Unable to create a new session")
		return errs.NewInternalError("%s", err)
	}

	c.Result = api.LoginResponse{
		SessionId:  *session.SessionId,
		ExpiryDate: session.ExpiryDate,
	}

	return nil
}

// FBLinkController Link the currently logged in user with the facebook user specified in the request
func FBLinkController(c *ctx.Context) errs.Error {
	var loginRequest api.FBLoginRequestData
	var err error
	if err = c.GinContext.BindJSON(&loginRequest); err != nil {
		return errs.NewRequestError("Request is invalid")
	}

	var fbUser *FBUser
	if fbUser, err = getFBUser(loginRequest.Token); err != nil {
		return errs.NewInternalError("Unable to get FB Identity")
	}

	var fbUserID = fbUser.Id
	var fbLink = fbUser.Link
	var userID = c.SessionData.UserId

	// link the user
	if err := linkFBUser(c.Db, userID, fbUserID, fbLink); err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

// linkFBUser link the specified user to the facebook user with fbUserID
func linkFBUser(db *gorm.DB, userID data.TUserID, fbUserID string, fbLink string) error {
	var externalAuthRecord *data.ExternalAuthData
	var err error
	if externalAuthRecord, err = query.GetExternalAuthRecord(db, userID); err != nil {
		return err
	}

	// update the user data
	if err := db.Model(&externalAuthRecord).
		Updates(&data.ExternalAuthData{FbUserId: &fbUserID, FbProfileLink: &fbLink}).Error; err != nil {
	}

	return nil
}

type FBUser struct {
	Id         string
	FirstName  string
	LastName   string
	Email      string
	Gender     data.GenderID
	Birthdate  string
	Link       string
	ProfilePic string
}

func getFBUser(accessToken string) (*FBUser, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "id,first_name,last_name,email,gender,birthday,link",
		"access_token": accessToken,
	})

	if err != nil {
		return nil, err
	}

	gender := data.GenderID(utility.GenderIdByName(res["gender"].(string)))

	if err != nil {
		return nil, errors.New("Unable to parse gender")
	}

	return &FBUser{
		Id:         res["id"].(string),
		FirstName:  res["first_name"].(string),
		LastName:   res["last_name"].(string),
		Email:      res["email"].(string),
		Gender:     gender,
		Birthdate:  res["birthday"].(string),
		Link:       res["link"].(string),
		ProfilePic: getFBProfilePicLink(res["id"].(string)),
	}, nil
}

func getFBProfilePicLink(userId string) string {
	return fmt.Sprintf(
		"http://graph.facebook.com/%s/picture?type=normal",
		userId,
	)
}
