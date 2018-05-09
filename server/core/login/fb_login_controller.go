package login

/**
 * Controller to handle user logging into
 */

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/getsentry/raven-go"
	fb "github.com/huandu/facebook"
	"github.com/romana/rlog"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/secrets"
	"letstalk/server/core/utility"
	"letstalk/server/data"
)

/**
 * Login with fb
 */
type FBLoginRequestData struct {
	Token             string `json:"token" binding:"required"`
	Expiry            int64  `json:"expiry" binding:"required"`
	NotificationToken string `json:"notificationToken"`
}

func FBController(c *ctx.Context) errs.Error {
	var loginRequest FBLoginRequestData
	var externalAuthRecord data.ExternalAuthData
	var userId int

	err := c.GinContext.BindJSON(&loginRequest)

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	authToken := loginRequest.Token
	expiry := time.Unix(loginRequest.Expiry, 0)

	user, err := getFBUser(authToken)
	db := c.Db

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	tx := c.Db.Begin()

	// check if the user already has facebook
	if tx.Where("fb_user_id = ?", user.Id).First(&externalAuthRecord).RecordNotFound() {

		appUser := data.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Gender:    user.Gender,
			Birthdate: user.Birthdate,
		}

		if err := tx.Where(&appUser).FirstOrCreate(&appUser).Error; err != nil {
			tx.Rollback()
			rlog.Error("Unable to insert new user")
			return errs.NewClientError("Unable to create user")
		}

		userId = appUser.UserId
		rlog.Debug("Created new user with id: ", userId)

		externalAuthRecord.FbUserId = &user.Id
		externalAuthRecord.UserId = userId
		// insert the user's fb auth data
		if err := tx.Create(&externalAuthRecord).Error; err != nil {
			rlog.Error(err)
			return errs.NewClientError("Unable to create user")
		}
		rlog.Debug("created auth record")

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
		// get a long lived access token from this short term token
		// do not fail if we cant do this
		go func() {
			rlog.Debug("Getting long lived fb token in another fiber.")

			res, err := fb.Get("/oauth/access_token", fb.Params{
				"grant_type":        "fb_exchange_token",
				"client_id":         secrets.GetSecrets().AppId,
				"client_secret":     secrets.GetSecrets().AppSecret,
				"fb_exchange_token": authToken,
			})

			if err != nil {
				// err can be an facebook API error.
				// if so, the Error struct contains error details.
				if e, ok := err.(*fb.Error); ok {
					rlog.Error("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]",
						e.Message, e.Type, e.Code, e.ErrorSubcode)
					raven.CaptureError(e, nil)
				}
			}

			if err != nil {
				// log the error to sentry
				// not fatal but this will cause early logout
				rlog.Error("Unable to get new token from facebook")
				raven.CaptureError(err, nil)
				return
			}
			expiresIn, err := res["expires_in"].(json.Number).Int64()
			if err != nil {
				rlog.Error("Malformed date")
				raven.CaptureError(err, nil)
				return
			}

			convertedTime := time.Unix(expiresIn, 0)
			if err != nil {
				rlog.Error("Unable to get new token from facebook")
				raven.CaptureError(err, nil)
				return
			}

			fbAuthToken := data.FbAuthToken{
				UserId:    userId,
				AuthToken: res["access_token"].(string),
				Expiry:    convertedTime,
			}

			if err := db.Save(&fbAuthToken).Error; err != nil {
				raven.CaptureError(err, nil)
				return
			}

			// TODO: we get a long term user access token but this might set of spam filter according to fb
			// https://developers.facebook.com/docs/facebook-login/access-tokens/expiration-and-extension
		}()
	} else {
		userId = externalAuthRecord.UserId
	}

	// create new session for user id
	session, err := (*c.SessionManager).CreateNewSessionForUserId(userId, &loginRequest.NotificationToken)

	if err != nil {
		rlog.Error("Unable to create a new session")
		return errs.NewInternalError("%s", err)
	}

	c.Result = LoginResponse{*session.SessionId, session.ExpiryDate}

	return nil
}

type FBUser struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Gender    int
	Birthdate *time.Time
}

func getFBUser(accessToken string) (*FBUser, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "id,first_name,last_name,email,gender,birthday",
		"access_token": accessToken,
	})

	if err != nil {
		return nil, err
	}

	gender := utility.GenderIdByName(res["gender"].(string))

	if err != nil {
		return nil, errors.New("Unable to parse gender")
	}

	birthdate, err := time.Parse("01/02/2006", res["birthday"].(string))
	if err != nil {
		return nil, errors.New("Unable to parse birthday")
	}

	return &FBUser{
		Id:        res["id"].(string),
		FirstName: res["first_name"].(string),
		LastName:  res["last_name"].(string),
		Email:     res["email"].(string),
		Gender:    gender,
		Birthdate: &birthdate,
	}, nil
}
