package routes

import (
	"letstalk/server/core/bootstrap"
	"letstalk/server/core/contact_info"
	"letstalk/server/core/controller"
	"letstalk/server/core/ctx"
	"letstalk/server/core/email_subscription"
	"letstalk/server/core/errs"
	"letstalk/server/core/login"
	"letstalk/server/core/matching"
	"letstalk/server/core/notifications"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/core/sessions"
	"net/http"
	"time"

	"letstalk/server/core/meeting"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type handlerWrapper struct {
	db *gorm.DB
	sm *sessions.ISessionManagerBase
}

type handlerFunc func(*ctx.Context) errs.Error

func Register(db *gorm.DB, sessionManager *sessions.ISessionManagerBase) *gin.Engine {
	hw := handlerWrapper{db, sessionManager}

	router := gin.Default()

	router.OPTIONS("/test")
	router.GET("/test", hw.wrapHandler(GetTest, false))

	router.OPTIONS("/testAuth")
	router.GET("/testAuth", hw.wrapHandler(GetTestAuth, true))

	v1 := router.Group("/v1")

	// create a new user
	v1.OPTIONS("/signup")
	v1.POST("/signup", hw.wrapHandler(login.SignupUser, false))

	// create a new session for an existing user
	v1.OPTIONS("/login")
	v1.POST("/login", hw.wrapHandler(login.LoginUser, false))

	// for fb_authentication
	v1.OPTIONS("/fb_login")
	v1.POST("/fb_login", hw.wrapHandler(login.FBController, false))

	v1.OPTIONS("/fb_link")
	v1.POST("/fb_link", hw.wrapHandler(login.FBLinkController, true))

	// update user data
	v1.OPTIONS("/cohort")
	v1.POST(
		"/cohort",
		hw.wrapHandler(onboarding.UpdateUserCohort, true),
	)
	v1.GET(
		"/cohort",
		hw.wrapHandler(query.GetCohortController, true),
	)

	// gets profile data about signed in user
	v1.OPTIONS("/me")
	v1.GET("/me", hw.wrapHandler(controller.GetMyProfileController, true))

	// updates profile data for signed in user
	v1.OPTIONS("/profile_edit")
	v1.POST("/profile_edit", hw.wrapHandler(controller.ProfileEditController, true))

	v1.OPTIONS("/contact_info")
	v1.GET("/contact_info", hw.wrapHandler(
		contact_info.GetContactInfoController,
		true),
	)

	v1.OPTIONS("/register_notification")
	v1.POST("/register_notification", hw.wrapHandler(
		notifications.GetNewNotificationToken,
		true),
	)

	v1.OPTIONS("/logout")
	v1.POST("/logout", hw.wrapHandler(
		login.LogoutHandler,
		true),
	)

	v1.OPTIONS("/user_vector")
	v1.POST("/user_vector", hw.wrapHandler(
		onboarding.UserVectorUpdateController,
		true,
	))

	// boostrap endpoints

	v1.OPTIONS("/bootstrap")
	v1.GET(
		"/bootstrap",
		hw.wrapHandler(bootstrap.GetCurrentUserBoostrapStatusController, true),
	)

	// request-to-match endpoints

	v1.OPTIONS("/all_credentials")
	v1.GET(
		"/all_credentials",
		hw.wrapHandler(controller.GetAllCredentialsController, false),
	)

	v1.OPTIONS("/credential")
	v1.POST(
		"/credential",
		hw.wrapHandler(controller.AddUserCredentialController, true),
	)
	v1.DELETE(
		"/credential",
		hw.wrapHandler(controller.RemoveUserCredentialController, true),
	)

	v1.OPTIONS("/credentials")
	v1.GET(
		"/credentials",
		hw.wrapHandler(controller.GetUserCredentialsController, true),
	)

	v1.OPTIONS("/credential_request")
	v1.POST(
		"/credential_request",
		hw.wrapHandler(controller.AddUserCredentialRequestController, true),
	)
	v1.DELETE(
		"/credential_request",
		hw.wrapHandler(controller.RemoveUserCredentialRequestController, true),
	)

	v1.OPTIONS("/credential_requests")
	v1.GET(
		"/credential_requests",
		hw.wrapHandler(controller.GetUserCredentialRequestsController, true),
	)

	v1.OPTIONS("/upload_profile_pic")
	v1.POST(
		"/upload_profile_pic",
		hw.wrapHandler(onboarding.ProfilePicController, true),
	)

	v1.OPTIONS("/subscribe_email")
	v1.POST(
		"/subscribe_email",
		hw.wrapHandler(email_subscription.AddSubscription, false),
	)

	// Meetings
	v1.OPTIONS("/meeting/confirm")
	v1.POST("/meeting/confirm", hw.wrapHandler(meeting.PostMeetingConfirmation, true /* auth required */))

	// Debug route group.
	debug := router.Group("/debug")

	debug.OPTIONS("/matching")
	debug.POST("/matching", hw.wrapHandler(matching.PostMatchingController, true /* auth required */))

	return router
}

/**
 * Wraps all requests.
 * If a header contains a sessionId attribute, we try to find an appropriate session
 */
func (hw handlerWrapper) wrapHandler(handler handlerFunc, needAuth bool) gin.HandlerFunc {
	return func(g *gin.Context) {
		var session *sessions.SessionData

		c := ctx.NewContext(g, hw.db, session, hw.sm)

		rlog.Debug("Checking if Auth needed")
		// the api route requires authentication so we have a session Id
		if needAuth {
			sessionId := g.GetHeader("sessionId")

			// check that the user provided a session id
			if sessionId == "" {
				rlog.Info("No session id provided.")
				c.GinContext.JSON(
					401,
					gin.H{"Error": query.Error{Code: 401, Message: "No session id provided. This is required."}},
				)
				return
			}

			session, err := (*hw.sm).GetSessionForSessionId(sessionId)

			// check that the session Id corresponds to an existing session
			if err != nil {
				rlog.Infof("%s", err)
				c.GinContext.JSON(401, gin.H{"Error": query.Error{Code: 401, Message: "Bad session Id."}})
				return
			}

			// check that the session token is not expired.
			if session.ExpiryDate.Before(time.Now()) {
				rlog.Error("Session token expired.")
				c.GinContext.JSON(401, gin.H{"Error": query.Error{Code: 401, Message: "Session token expired."}})
				return
			}
			c.SessionData = session

		}

		rlog.Debug("Running handler")

		err := handler(c)

		if err != nil {
			rlog.Infof("Returning error: %s\n", err)
			raven.CaptureError(err, nil)
			c.GinContext.JSON(err.GetHTTPCode(), gin.H{"Error": convertError(err)})
			return
		}
		rlog.Infof("Returning result: %v\n", c.Result)
		c.GinContext.JSON(http.StatusOK, gin.H{"Result": c.Result})
	}
}

func convertError(e errs.Error) query.Error {
	return query.Error{
		Code:    e.GetHTTPCode(),
		Message: e.Error(),
	}
}
