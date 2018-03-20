package login

/**
 * Controller to handle user logging into
 */

import (
	"github.com/getsentry/raven-go"
	fb "github.com/huandu/facebook"
	"letstalk/server/core/ctx"
	"letstalk/server/core/db"
	"letstalk/server/core/errs"
	"letstalk/server/core/secrets"
)

/**
 * Login with fb
 */
type FBLoginRequestData struct {
	Token             string
	NotificationToken string
}

func FBController(c *ctx.Context) errs.Error {
	var loginRequest FBLoginRequestData
	err := c.GinContext.BindJSON(loginRequest)

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	authToken := loginRequest.Token

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
			return errs.NewDbError(err)
		}

		// get a unique id
		// TODO chagne db.NumId interface to not take a context and to use a transaction
		if userId, err = db.NumId(c); err != nil {
			tx.Rollback()
			return errs.NewDbError(err)
		}

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
			errs.NewDbError(err)
		}

		// initially insert with blank expiry indicating that
		// this is a short lived token
		_, err = fb_auth_token_stmt.Exec(userId, authToken, "")

		if err != nil {
			tx.Rollback()
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
		_, err = fb_mapping_stmt.Exec(user.Id)

		if err != nil {
			tx.Rollback()
			return errs.NewDbError(err)
		}

		err = tx.Commit()
		if err != nil {
			return errs.NewDbError(err)
		}

		// get a long lived access token from this short term token
		// do not fail if we cant do this
		go func() {
			res, err := fb.Get("/oauth/access_token", fb.Params{
				"grant_type":        "fb_exchange_token",
				"client_id":         secrets.GetSecrets().AppId,
				"client_secret":     secrets.GetSecrets().AppSecret,
				"fb_exchange_token": authToken,
			})

			if err != nil {
				// log the error to sentry
				// not fatal but this will cause early logout
				raven.CaptureError(err, nil)
			}
			// insert the fb auth data
			fb_auth_token_stmt, err := c.Db.Prepare(
				`
				INSERT INTO fb_auth_token
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
				raven.CaptureError(err, nil)
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
	Gender    string
	Birthdate string
}

func getFBUser(accessToken string) (*FBUser, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "id,first_name,last_name,email,gender,birthday",
		"access_token": accessToken,
	})

	if err != nil {
		return nil, err
	}

	return &FBUser{
		Id:        res["id"].(string),
		FirstName: res["first_name"].(string),
		LastName:  res["last_name"].(string),
		Email:     res["email"].(string),
		Gender:    res["gender"].(string),
		Birthdate: res["birthday"].(string),
	}, nil
}
