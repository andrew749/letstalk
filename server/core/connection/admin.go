package connection

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/core/notifications"
)

/**
 * AddMentorshipController is an admin function that adds a new mentorship connection.
 */
func AddMentorshipController(c *ctx.Context) errs.Error {
	var input api.CreateMentorship
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if err := handleAddMentorship(c, &input); err != nil {
		return err
	}
	if err := sendMentorshipNotifications(c, &input); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func handleAddMentorship(c *ctx.Context, request *api.CreateMentorship) errs.Error {
	if request.MentorId == request.MenteeId {
		return errs.NewRequestError("mentor and mentee user must be different")
	}
	if user, err := query.GetUserById(c.Db, request.MentorId); err != nil || user == nil {
		return errs.NewRequestError("no such user %d", request.MentorId)
	}
	if user, err := query.GetUserById(c.Db, request.MenteeId); err != nil || user == nil {
		return errs.NewRequestError("no such user %d", request.MenteeId)
	}
	if conn, err := query.GetConnectionDetailsUndirected(c.Db, request.MentorId, request.MenteeId); err != nil {
		return errs.NewDbError(err)
	} else if conn != nil {
		return errs.NewRequestError("connection already exists")
	}
	intent := data.ConnectionIntent{
		Type: data.INTENT_TYPE_ASSIGNED,
	}
	createdAt := time.Now()
	mentorship := data.Mentorship{
		MentorUserId: request.MentorId,
		CreatedAt: createdAt,
	}
	conn := data.Connection{
		UserOneId: request.MentorId,
		UserTwoId: request.MenteeId,
		CreatedAt: createdAt,
		AcceptedAt: &createdAt, // Automatically accept.
		Intent: &intent,
		Mentorship: &mentorship,
	}
	if err := c.Db.Create(&conn).Error; err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func sendMentorshipNotifications(c *ctx.Context, request *api.CreateMentorship) errs.Error {
	// Send notifications to matched pair.
	notifErr1 := notifications.NewMenteeNotification(c.Db, request.MentorId, request.MenteeId)
	notifErr2 := notifications.NewMentorNotification(c.Db, request.MenteeId, request.MentorId)
	if notifErr1 != nil || notifErr2 != nil {
		return errs.NewInternalError("error sending user notifications: %v; %v", notifErr1, notifErr2)
	}
	return nil
}

