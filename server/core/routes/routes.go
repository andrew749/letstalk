package routes

import (
	"letstalk/server/core/auth"
	"letstalk/server/core/bootstrap"
	"letstalk/server/core/contact_info"
	"letstalk/server/core/controller"
	"letstalk/server/core/ctx"
	"letstalk/server/core/email_subscription"
	"letstalk/server/core/errs"
	"letstalk/server/core/matching"
	"letstalk/server/core/meeting"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/core/sessions"
	"letstalk/server/core/user"
	"net/http"
	"time"

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

func debugAuthMiddleware(db *gorm.DB, sessionManager *sessions.ISessionManagerBase) gin.HandlerFunc {
	return func(g *gin.Context) {
		session, err := getSessionData(g, sessionManager)
		if err != nil {
			abortWithError(g, err)
			return
		}
		authUser, e := query.GetUserById(db, session.UserId)
		if e != nil || !auth.HasAdminAccess(authUser) {
			g.AbortWithStatusJSON(http.StatusNotFound, "404 page not found")
			return
		}
		g.Next()
	}
}

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
	v1.POST("/signup", hw.wrapHandler(user.SignupUser, false))

	// create a new session for an existing user
	v1.OPTIONS("/login")
	v1.POST("/login", hw.wrapHandler(user.LoginUser, false))

	// user forgot password, generate a link to send to email
	v1.OPTIONS("/forgot_password")
	v1.POST("/forgot_password", hw.wrapHandler(user.GenerateNewForgotPasswordRequestController, false))

	// handle changin a user password using unique link.
	v1.OPTIONS("/change_password")
	v1.POST("/change_password", hw.wrapHandler(user.ForgotPasswordController, false))

	// for fb_authentication
	v1.OPTIONS("/fb_login")
	v1.POST("/fb_login", hw.wrapHandler(user.FBController, false))

	v1.OPTIONS("/fb_link")
	v1.POST("/fb_link", hw.wrapHandler(user.FBLinkController, true))

	// update user data
	v1.OPTIONS("/cohort")
	v1.POST(
		"/cohort",
		hw.wrapHandler(onboarding.UpdateUserCohort, true),
	)
	v1.GET(
		"/cohort",
		hw.wrapHandler(controller.GetCohortController, true),
	)

	// update user data
	v1.OPTIONS("/cohorts")
	v1.GET(
		"/cohorts",
		hw.wrapHandler(controller.GetAllCohortsController, false),
	)

	// gets profile data about signed in user
	v1.OPTIONS("/me")
	v1.GET("/me", hw.wrapHandler(user.GetMyProfileController, true))

	// gets profile data about a match for signed in user
	v1.OPTIONS("/match_profile/:userId")
	v1.GET("/match_profile/:userId", hw.wrapHandler(user.GetMatchProfileController, true))

	// gets profile data about a match for signed in user
	v1.OPTIONS("/remove_rtm_matches/:userId")
	v1.DELETE("/remove_rtm_matches/:userId", hw.wrapHandler(controller.RemoveRtmMatches, true))

	v1.OPTIONS("/profile_pic")
	v1.GET("/profile_pic", hw.wrapHandler(user.GetProfilePicUrl, true))

	// updates profile data for signed in user
	v1.OPTIONS("/profile_edit")
	v1.POST("/profile_edit", hw.wrapHandler(user.ProfileEditController, true))

	v1.OPTIONS("/contact_info")
	v1.GET("/contact_info", hw.wrapHandler(
		contact_info.GetContactInfoController,
		true),
	)

	v1.OPTIONS("/register_notification")
	v1.POST("/register_notification", hw.wrapHandler(
		controller.GetNewNotificationToken,
		true),
	)

	v1.OPTIONS("/logout")
	v1.POST("/logout", hw.wrapHandler(
		user.LogoutHandler,
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

	// Notifications
	v1.OPTIONS("/notifications")
	v1.GET("/notifications", hw.wrapHandler(controller.GetNotifications, true))

	v1.OPTIONS("/notifications/update_state")
	v1.POST("/notifications/update_state", hw.wrapHandler(controller.UpdateNotificationState, true))

	// Debug route group.
	debug := router.Group("/debug")
	debug.Use(debugAuthMiddleware(hw.db, hw.sm))

	debug.OPTIONS("/matching")
	debug.POST("/matching", hw.wrapHandler(matching.PostMatchingController, false))

	return router
}

/**
 * Wraps all requests.
 * If a header contains a sessionId attribute, we try to find an appropriate session
 */
func (hw handlerWrapper) wrapHandler(handler handlerFunc, needAuth bool) gin.HandlerFunc {
	return func(g *gin.Context) {
		c := ctx.NewContext(g, hw.db, nil /* session */, hw.sm)

		// The api route requires authentication so we add session data from the header.
		if needAuth {
			session, err := getSessionData(g, hw.sm)
			if err != nil {
				abortWithError(c.GinContext, err)
				return
			}
			c.SessionData = session
		}

		rlog.Debug("Running handler")
		err := handler(c)

		if err != nil {
			abortWithError(c.GinContext, err)
			return
		}

		rlog.Infof("Returning result: %v\n", c.Result)
		c.GinContext.JSON(http.StatusOK, gin.H{"Result": c.Result})
	}
}

func abortWithError(g *gin.Context, err errs.Error) {
	rlog.Errorf("Returning error: %s\n", err)
	raven.CaptureError(err, nil)
	g.AbortWithStatusJSON(err.GetHTTPCode(), gin.H{"Error": convertError(err)})
}

func convertError(e errs.Error) query.Error {
	return query.Error{
		Code:    e.GetHTTPCode(),
		Message: e.Error(),
	}
}

func getSessionData(g *gin.Context, sessionManager *sessions.ISessionManagerBase) (*sessions.SessionData, errs.Error) {
	sessionId := g.GetHeader("sessionId")

	// check that the user provided a session id
	if sessionId == "" {
		rlog.Info("No session id provided.")
		return nil, errs.NewUnauthorizedError("Required session id not provided")
	}

	session, err := (*sessionManager).GetSessionForSessionId(sessionId)

	// check that the session Id corresponds to an existing session
	if err != nil {
		rlog.Infof("%s", err)
		return nil, errs.NewUnauthorizedError("Bad session id")
	}

	// check that the session token is not expired.
	if session.ExpiryDate.Before(time.Now()) {
		rlog.Error("Session token expired.")
		return nil, errs.NewUnauthorizedError("Session token expired.")
	}
	return session, nil
}
