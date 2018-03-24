package login

/**
 * Controller to handle user logging into
 */

import (
	"errors"
	"time"

	"github.com/getsentry/raven-go"
	fb "github.com/huandu/facebook"
	"github.com/romana/rlog"
	"letstalk/server/core/ctx"
	"letstalk/server/core/db"
	"letstalk/server/core/errs"
	"letstalk/server/core/secrets"
	"letstalk/server/core/utility"
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
	err := c.GinContext.BindJSON(&loginRequest)

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	authToken := loginRequest.Token
	expiry := time.Unix(loginRequest.Expiry, 0)

	user, err := getFBUser(authToken)

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	tx, err := c.Db.Beginx()

	if err != nil {
		return errs.NewDbError(err)
	}

	// check if the user id exist in our mappings
	stmt, err := tx.Prepare(
		`
		SELECT fb_auth_data.user_id
		FROM fb_auth_data
		INNER JOIN user ON user.user_id=fb_auth_data.user_id
		WHERE fb_auth_data.fb_user_id = ?
		`,
	)

	if err != nil {
		return errs.NewDbError(err)
	}

	rows, err := stmt.Query(user.Id)

	if err != nil {
		return errs.NewDbError(err)
	}

	// the user Id for this application
	var userId int

	// if the user doesn't have an account
	if !rows.Next() {
		stmt, err = tx.Prepare(
			`
			INSERT INTO user
				(user_id, first_name, last_name, email, gender, birthdate)
			VALUES (?, ?, ?, ?, ?, ?)
			`,
		)

		if err != nil {
			tx.Rollback()
			rlog.Error("Unable to prepare db.", userId)
			return errs.NewDbError(err)
		}

		// get a unique id
		// TODO change db.NumId interface to not take a context and to use a transaction instead
		if userId, err = db.NumId(c); err != nil {
			tx.Rollback()
			rlog.Error("Unable to generate id.", userId)
			return errs.NewDbError(err)
		}
		rlog.Info("Registering user with id: ", userId)

		_, err = stmt.Exec(
			userId,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Gender,
			user.Birthdate,
		)

		// if there was an issue inserting this user
		if err != nil {
			tx.Rollback()
			rlog.Error("Unable to insert new user")
			return errs.NewDbError(err)
		}

		// insert the fb auth data
		fb_auth_token_stmt, err := tx.Prepare(
			`
			INSERT INTO fb_auth_token
				(user_id, auth_token, expiry)
			VALUES (?, ?, ?)
			`,
		)
		if err != nil {
			tx.Rollback()
			rlog.Error("Unable to prepare auth token")
			errs.NewDbError(err)
		}

		// initially insert with blank expiry indicating that
		// this is a short lived token
		_, err = fb_auth_token_stmt.Exec(userId, authToken, expiry)

		if err != nil {
			tx.Rollback()
			rlog.Error("Unable to insert auth token")
			return errs.NewDbError(err)
		}

		fb_mapping_stmt, err := tx.Prepare(
			`
			INSERT INTO fb_auth_data
				(user_id, fb_user_id)
			VALUES(?, ?)
			`,
		)

		// insert a new user mapping
		_, err = fb_mapping_stmt.Exec(userId, user.Id)

		if err != nil {
			tx.Rollback()
			rlog.Error("Unable to insert facebook token")
			return errs.NewDbError(err)
		}

		err = tx.Commit()
		if err != nil {
			rlog.Error("Unable to commit everything")
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

			// insert the non-stale fb auth data
			fb_auth_token_stmt, err := c.Db.Prepare(
				`
				REPLACE INTO fb_auth_token
					(user_id, auth_token, expiry)
				VALUES (?, ?, ?)
				`,
			)

			_, err = fb_auth_token_stmt.Exec(
				userId,
				res["access_token"].(string),
				res["expires_in"].(string),
			)

			if err != nil {
				rlog.Error("Unable to insert newer token in db.")
				raven.CaptureError(err, nil)
				return
			}

			// TODO: we get a long term user access token but this might set of spam filter according to fb
			// https://developers.facebook.com/docs/facebook-login/access-tokens/expiration-and-extension
		}()
	} else {
		err := rows.Scan(&userId)
		if err != nil {
			return errs.NewDbError(err)
		}
	}

	sm := c.SessionManager
	// create new session for user id
	session, err := (*sm).CreateNewSessionForUserId(userId, &loginRequest.NotificationToken)

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
