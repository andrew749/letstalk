package connection

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

/**
 * AddMentorshipController is an admin function that adds a new mentorship connection.
 */
func AddMentorshipController(c *ctx.Context) errs.Error {
	var input api.CreateMentorshipByEmail
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if err := handleAddMentorship(c.Db, &input); err != nil {
		return err
	}
	if err := sendMentorshipNotifications(c.Db, &input); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func handleAddMentorship(db *gorm.DB, request *api.CreateMentorshipByEmail) errs.Error {
	var mentor, mentee *data.User
	var err error
	if request.MentorEmail == request.MenteeEmail {
		return errs.NewRequestError("mentor and mentee user must be different")
	}
	if mentor, err = query.GetUserByEmail(db, request.MentorEmail); err != nil || mentor == nil {
		return errs.NewRequestError("no such user %s", request.MentorEmail)
	}
	if mentee, err = query.GetUserByEmail(db, request.MenteeEmail); err != nil || mentee == nil {
		return errs.NewRequestError("no such user %s", request.MenteeEmail)
	}
	if conn, err := query.GetConnectionDetailsUndirected(db, mentor.UserId, mentee.UserId); err != nil {
		return errs.NewDbError(err)
	} else if conn != nil {
		return errs.NewRequestError("connection already exists")
	}
	intent := data.ConnectionIntent{
		Type: data.INTENT_TYPE_ASSIGNED,
	}
	createdAt := time.Now()
	mentorship := data.Mentorship{
		MentorUserId: mentor.UserId,
		CreatedAt: createdAt,
	}
	conn := data.Connection{
		UserOneId: mentor.UserId,
		UserTwoId: mentee.UserId,
		CreatedAt: createdAt,
		AcceptedAt: &createdAt, // Automatically accept.
		Intent: &intent,
		Mentorship: &mentorship,
	}
	if err := db.Create(&conn).Error; err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func sendMentorshipNotifications(db *gorm.DB, request *api.CreateMentorshipByEmail) errs.Error {
	mentor, err := query.GetUserByEmail(db, request.MentorEmail)
	if err != nil {
		return errs.NewDbError(err)
	}
	mentee, _ := query.GetUserByEmail(db, request.MenteeEmail)
	if err != nil {
		return errs.NewDbError(err)
	}
	// Send notifications to matched pair.
	notifErr1 := notifications.NewMenteeNotification(db, mentor.UserId, mentee.UserId)
	notifErr2 := notifications.NewMentorNotification(db, mentee.UserId, mentor.UserId)
	if notifErr1 != nil || notifErr2 != nil {
		return errs.NewInternalError("error sending user notifications: %v; %v", notifErr1, notifErr2)
	}
	return nil
}
