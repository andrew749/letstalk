package connection

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/email"

	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"letstalk/server/core/meetup_reminder"
)

/**
 * AddMentorshipController is an admin function that adds a new mentorship connection.
 */
func AddMentorshipController(c *ctx.Context) errs.Error {
	var input api.CreateMentorshipByEmail
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if err := HandleAddMentorship(c.Db, &input); err != nil {
		rlog.Error("failed to add mentorship for mentor/mentee pair", input)
		return err
	}
	if input.RequestType == api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN {
		if err := sendMentorshipNotifications(c.Db, &input); err != nil {
			return err
		}
	}
	c.Result = "Ok"
	return nil
}

// TODO(wojtechnology): Give this a more explicit public interface.
func HandleAddMentorship(db *gorm.DB, request *api.CreateMentorshipByEmail) errs.Error {
	var mentor, mentee *data.User
	var err errs.Error
	if request.MentorEmail == request.MenteeEmail {
		return errs.NewRequestError("mentor and mentee user must be different")
	}
	var noSuchMentorErr, noSuchMenteeErr string
	if mentor, err = query.GetUserByEmail(db, request.MentorEmail); err != nil || mentor == nil {
		noSuchMentorErr = fmt.Sprintf("no such user %s, %+v", request.MentorEmail, err)
	}
	if mentee, err = query.GetUserByEmail(db, request.MenteeEmail); err != nil || mentee == nil {
		noSuchMenteeErr = fmt.Sprintf("no such user %s, %+v", request.MenteeEmail, err)
	}
	if len(noSuchMentorErr) > 0 || len(noSuchMenteeErr) > 0 {
		return errs.NewNotFoundError("%s %s", noSuchMentorErr, noSuchMenteeErr)
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
		CreatedAt:    createdAt,
	}
	conn := data.Connection{
		UserOneId:  mentor.UserId,
		UserTwoId:  mentee.UserId,
		CreatedAt:  createdAt,
		AcceptedAt: &createdAt, // Automatically accept.
		Intent:     &intent,
		Mentorship: &mentorship,
	}
	if request.RequestType == api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN {
		rlog.Infof("not a dry run, adding (%s, %s)", mentor.Email, mentee.Email)
		dbErr := ctx.WithinTx(db, func(tx *gorm.DB) error {
			if err := tx.Create(&conn).Error; err != nil {
				return err
			}
			if err := meetup_reminder.ScheduleInitialReminder(tx, conn.UserOneId, conn.UserTwoId); err != nil {
				return err
			}
			return nil
		})
		if dbErr != nil {
			return errs.NewDbError(err)
		}
	}
	return nil
}

func sendMentorshipNotifications(db *gorm.DB, request *api.CreateMentorshipByEmail) errs.Error {
	mentor, err := query.GetUserByEmail(db, request.MentorEmail)
	if err != nil {
		return err
	}

	mentee, err := query.GetUserByEmail(db, request.MenteeEmail)
	if err != nil {
		return err
	}

	mentorCohort, err := query.GetUserCohort(db, mentor.UserId)
	if err != nil {
		return err
	}

	menteeCohort, err := query.GetUserCohort(db, mentee.UserId)
	if err != nil {
		return err
	}

	// Send notifications to matched pair.
	notifErr1 := notifications.NewMenteeNotification(db, mentor.UserId, mentee)
	notifErr2 := notifications.NewMentorNotification(db, mentee.UserId, mentor)

	mentorEmail := mail.NewEmail(mentor.FirstName, mentor.Email)
	menteeEmail := mail.NewEmail(mentee.FirstName, mentee.Email)

	emailErr1 := email.SendNewMenteeEmail(mentorEmail, mentor.FirstName, mentee.FirstName, menteeCohort.ProgramName, menteeCohort.GradYear)
	emailErr2 := email.SendNewMentorEmail(menteeEmail, mentor.FirstName, mentee.FirstName, mentorCohort.ProgramName, mentorCohort.GradYear)
	var compositeError *errs.CompositeError
	compositeError = errs.AppendNullableError(compositeError, notifErr1)
	compositeError = errs.AppendNullableError(compositeError, notifErr2)
	compositeError = errs.AppendNullableError(compositeError, emailErr1)
	compositeError = errs.AppendNullableError(compositeError, emailErr2)
	if emailErr1 != nil {
		rlog.Errorf("Unable to send email to user %d with email: %s;Error: %+v", mentor.UserId, mentor.Email, emailErr1)
	}

	if emailErr2 != nil {
		rlog.Errorf("Unable to send email to user %d with email: %s;Error: %+v", mentee.UserId, mentee.Email, emailErr2)
	}

	if notifErr1 != nil {
		rlog.Errorf("Unable to send notification to user %d with email: %s;Error: %+v", mentor.UserId, mentor.Email, notifErr1)
	}

	if notifErr2 != nil {
		rlog.Errorf("Unable to send notification to user %d with email: %s;Error: %+v", mentee.UserId, mentee.Email, notifErr2)
	}

	if compositeError != nil {
		rlog.Errorf("%+v", compositeError)
		return compositeError
	}
	return nil
}
