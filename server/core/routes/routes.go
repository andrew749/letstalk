package routes

import (
	"net/http"
	"time"

	"letstalk/server/core/auth"
	"letstalk/server/core/bootstrap"
	"letstalk/server/core/connection"
	"letstalk/server/core/contact_info"
	"letstalk/server/core/controller"
	"letstalk/server/core/ctx"
	"letstalk/server/core/email_subscription"
	"letstalk/server/core/errs"
	"letstalk/server/core/match_round"
	"letstalk/server/core/matching"
	"letstalk/server/core/meeting"
	"letstalk/server/core/meetup_reminder"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/core/sessions"
	"letstalk/server/core/user"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

type handlerWrapper struct {
	db *gorm.DB
	es *elastic.Client
	sm *sessions.ISessionManagerBase
}

type handlerFunc func(*ctx.Context) errs.Error

func adminAuthMiddleware(db *gorm.DB, sessionManager *sessions.ISessionManagerBase) gin.HandlerFunc {
	return func(g *gin.Context) {
		session, err := getSessionData(g, sessionManager)
		if err != nil {
			abortWithError(g, err)
			return
		}
		authUser, e := query.GetUserById(db, session.UserId)
		if e != nil || !auth.HasAdminAccess(authUser) {
			rlog.Infof("rejected non-admin user with id %d", authUser.UserId)
			g.AbortWithStatusJSON(http.StatusNotFound, "404 page not found")
			return
		}
		g.Next()
	}
}

func Register(
	db *gorm.DB,
	es *elastic.Client,
	sessionManager *sessions.ISessionManagerBase,
) *gin.Engine {
	hw := handlerWrapper{db, es, sessionManager}

	router := gin.Default()

	router.OPTIONS("/test")
	router.GET("/test", hw.wrapHandler(GetTest, false))

	router.OPTIONS("/testAuth")
	router.GET("/testAuth", hw.wrapHandler(GetTestAuth, true))

	router.LoadHTMLGlob("templates/*")
	router.LoadHTMLGlob("web/dist/*.html")
	router.Static("/assets", "web/dist/assets/")

	// Html login page
	router.OPTIONS("/admin_panel")
	router.GET("/admin_panel/*any", hw.wrapHandlerHTML(controller.GetAdminPanel, false))

	router.OPTIONS("/web")
	router.GET("/web/*any", hw.wrapHandlerHTML(controller.GetWebapp, false))

	// Render a page to make it easy to send notification campaigns
	router.OPTIONS("/notification_console")
	router.GET("/notification_console", hw.wrapHandlerHTML(controller.GetNotificationManagementConsole, false))

	v1 := router.Group("/v1")

	// additional asset routes
	v1.Static("/assets", "web/dist/assets/")

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

	// endpoint to send a new account verification email
	v1.OPTIONS("/send_email_verification")
	v1.POST("/send_email_verification", hw.wrapHandler(user.SendEmailVerificationController, true))

	// callback to verify user's email
	v1.OPTIONS("/verify_email")
	v1.POST("/verify_email", hw.wrapHandler(user.VerifyEmailController, false))

	// callback to verify a link in the user_verify_link table
	v1.OPTIONS("/verify_link")
	v1.POST("/verify_link", hw.wrapHandler(controller.VerifyLinkController, false))

	// for fb_authentication
	v1.OPTIONS("/fb_login")
	v1.POST("/fb_login", hw.wrapHandler(user.FBController, false))

	v1.OPTIONS("/fb_link")
	v1.POST("/fb_link", hw.wrapHandler(user.FBLinkController, true))

	// update user data
	v1.OPTIONS("/cohort")
	v1.POST(
		"/cohort",
		hw.wrapHandler(controller.UpdateUserCohortAndAdditionalInfo, true),
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

	// gets profile data about a scanned user by their qr code
	v1.OPTIONS("/public_profile/:code")
	v1.GET("/public_profile/:code", hw.wrapHandler(user.GetPublicProfileController, true))

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

	v1.OPTIONS("/logout")
	v1.POST("/logout", hw.wrapHandler(
		user.LogoutHandler,
		true),
	)

	// boostrap endpoints

	v1.OPTIONS("/bootstrap")
	v1.GET(
		"/bootstrap",
		hw.wrapHandler(bootstrap.GetCurrentUserBoostrapStatusController, true),
	)

	// request-to-match endpoints

	v1.OPTIONS("/connection")
	// request a new connection with another user
	v1.POST("/connection", hw.wrapHandler(connection.PostRequestConnection, true))
	// remove a connection
	v1.DELETE("/connection", hw.wrapHandler(connection.RemoveConnection, true))

	// accept a connection request from another user
	v1.OPTIONS("/connection/accept")
	v1.POST("/connection/accept", hw.wrapHandler(connection.PostAcceptConnection, true))

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
		hw.wrapHandler(controller.UploadProfilePic, true),
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

	v1.OPTIONS("/notification")
	v1.GET("/notification/:notificationId", hw.wrapHandler(controller.GetNotification, true))

	v1.OPTIONS("/notifications/update_state")
	v1.POST("/notifications/update_state", hw.wrapHandler(controller.UpdateNotificationState, true))

	v1.OPTIONS("/notification_page")
	v1.GET("/notification_page", hw.wrapHandlerHTML(notifications.GetNotificationContentPage, true))

	v1.OPTIONS("/echo_notification")
	v1.POST("/echo_notification", hw.wrapHandlerHTML(notifications.EchoNotificationPage, true))

	// User Simple Traits
	v1.OPTIONS("/user_simple_trait")
	v1.POST("/user_simple_trait", hw.wrapHandler(controller.AddUserSimpleTraitByIdController, true))
	v1.DELETE(
		"/user_simple_trait",
		hw.wrapHandler(controller.RemoveUserSimpleTraitController, true),
	)

	v1.OPTIONS("/user_simple_trait_by_name")
	v1.POST(
		"/user_simple_trait_by_name",
		hw.wrapHandler(controller.AddUserSimpleTraitByNameController, true),
	)

	// User Devices
	v1.OPTIONS("/user_device/expo")
	v1.POST("/user_device/expo", hw.wrapHandler(controller.AddExpoDeviceToken, true))

	// User Positions
	v1.OPTIONS("/user_position")
	v1.POST("/user_position", hw.wrapHandler(controller.AddUserPositionController, true))
	v1.DELETE(
		"/user_position",
		hw.wrapHandler(controller.RemoveUserPositionController, true),
	)

	// User Group
	v1.OPTIONS("/user_group")
	v1.POST("/user_group", hw.wrapHandler(controller.AddUserGroupController, true))
	v1.DELETE(
		"/user_group",
		hw.wrapHandler(controller.RemoveUserGroupController, true),
	)

	// User surveys
	v1.OPTIONS("/survey")
	v1.POST("/survey", hw.wrapHandler(controller.PostSurveyResponses, true /* auth required */))
	v1.OPTIONS("/survey/:group")
	v1.GET("/survey/:group", hw.wrapHandler(controller.GetSurvey, true /* auth required */))

	// Meetup reminders
	v1.OPTIONS("/meetup_reminder")
	v1.POST("/meetup_reminder", hw.wrapHandler(meetup_reminder.PostMeetupReminder, true /* auth required */))
	v1.DELETE("/meetup_reminder", hw.wrapHandler(meetup_reminder.DeleteMeetupReminder, true /* auth required */))

	// Autocomplete endpoints
	autocompleteV1 := v1.Group("/autocomplete")
	autocompleteV1.OPTIONS("/simple_trait")
	autocompleteV1.POST(
		"/simple_trait",
		hw.wrapHandler(controller.SimpleTraitAutocompleteController, false),
	)

	autocompleteV1.OPTIONS("/role")
	autocompleteV1.POST(
		"/role",
		hw.wrapHandler(controller.RoleAutocompleteController, false),
	)

	autocompleteV1.OPTIONS("/organization")
	autocompleteV1.POST(
		"/organization",
		hw.wrapHandler(controller.OrganizationAutocompleteController, false),
	)

	autocompleteV1.OPTIONS("/multi_trait")
	autocompleteV1.POST(
		"/multi_trait",
		hw.wrapHandler(controller.MultiTraitAutocompleteController, false),
	)

	// User search endpoints
	userSearchV1 := v1.Group("/user_search")

	userSearchV1.OPTIONS("/simple_trait")
	userSearchV1.POST("/simple_trait", hw.wrapHandler(controller.SimpleTraitUserSearchController, true))

	userSearchV1.OPTIONS("/cohort")
	userSearchV1.POST("/cohort", hw.wrapHandler(controller.CohortUserSearchController, true))

	userSearchV1.OPTIONS("/my_cohort")
	userSearchV1.POST("/my_cohort", hw.wrapHandler(controller.MyCohortUserSearchController, true))

	userSearchV1.OPTIONS("/position")
	userSearchV1.POST("/position", hw.wrapHandler(controller.PositionUserSearchController, true))

	userSearchV1.OPTIONS("/group")
	userSearchV1.POST("/group", hw.wrapHandler(controller.GroupUserSearchController, true))

	// Admin route group.
	admin := router.Group("/admin")
	admin.Use(adminAuthMiddleware(hw.db, hw.sm))

	admin.OPTIONS("/matching")
	admin.POST("/matching", hw.wrapHandler(matching.PostMatchingController, false))

	admin.OPTIONS("/mentorship")
	admin.POST("/mentorship", hw.wrapHandler(connection.AddMentorshipController, false))

	admin.OPTIONS("/adhoc_notification")
	admin.POST("/adhoc_notification", hw.wrapHandler(controller.SendAdhocNotification, false))

	admin.OPTIONS("/campaign")
	admin.POST("/campaign", hw.wrapHandler(controller.NotificationCampaignController, false))

	admin.OPTIONS("/nuke_user")
	admin.POST("/nuke_user", hw.wrapHandler(controller.NukeUser, false))

	admin.OPTIONS("/create_match_round")
	admin.POST("/create_match_round", hw.wrapHandler(match_round.CreateMatchRoundController, false))

	admin.OPTIONS("/commit_match_round")
	admin.POST("/commit_match_round", hw.wrapHandler(match_round.CommitMatchRoundController, false))

	admin.OPTIONS("/match_rounds")
	admin.GET("/match_rounds", hw.wrapHandler(match_round.GetMatchRoundsController, false))

	admin.OPTIONS("/match_round")
	admin.DELETE("/match_round", hw.wrapHandler(match_round.DeleteMatchRoundController, false))

	return router
}

/**
 * Wraps all requests.
 * If a header contains a sessionId attribute, we try to find an appropriate session
 */
func (hw handlerWrapper) wrapHandler(handler handlerFunc, needAuth bool) gin.HandlerFunc {
	return func(g *gin.Context) {
		c := ctx.NewContext(g, hw.db, hw.es, nil /* session */, hw.sm)

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

func (hw handlerWrapper) wrapHandlerHTML(handler handlerFunc, needAuth bool) gin.HandlerFunc {
	return func(g *gin.Context) {
		c := ctx.NewContext(g, hw.db, hw.es, nil /* session */, hw.sm)

		// The api route requires authentication so we add session data from the header.
		if needAuth {
			session, err := getSessionData(g, hw.sm)
			if err != nil {
				abortWithErrorHTML(c.GinContext, err)
				return
			}
			c.SessionData = session
		}

		rlog.Debug("Running handler")
		err := handler(c)

		if err != nil {
			abortWithErrorHTML(c.GinContext, err)
			return
		}

		rlog.Infof("Returning result: %v\n", c.Result)
	}
}

func abortWithErrorHTML(g *gin.Context, err errs.Error) {
	rlog.Errorf("Returning error: %s\n", err.VerboseError())
	raven.CaptureError(err, nil)
	g.AbortWithStatus(err.GetHTTPCode())
}

func abortWithError(g *gin.Context, err errs.Error) {
	rlog.Errorf("Returning error: %s\n", err.VerboseError())
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

	// if no session id, try checking Cookie
	if sessionId == "" {
		sessionId, _ = g.Cookie("sessionId")
	}

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
