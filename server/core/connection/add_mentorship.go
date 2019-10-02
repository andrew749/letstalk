package connection

import (
	"fmt"
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/meetup_reminder"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/email"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
		if err := sendMentorshipNotificationsHandler(c.Db, &input); err != nil {
			return err
		}
	}
	c.Result = "Ok"
	return nil
}

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
	tx := db.Begin()
	if err := AddMentorship(db, mentor.UserId, mentee.UserId, request.RequestType); err != nil {
		tx.Rollback()
		return err
	}
	if dbErr := tx.Commit().Error; dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	return nil
}

func AddMentorship(
	tx *gorm.DB,
	mentorUserId data.TUserID,
	menteeUserId data.TUserID,
	requestType api.CreateMentorshipType,
) errs.Error {
	if mentorUserId == menteeUserId {
		return errs.NewRequestError("mentor and mentee user must be different")
	}
	if conn, err := query.GetConnectionDetailsUndirected(tx, mentorUserId, menteeUserId); err != nil {
		return errs.NewDbError(err)
	} else if conn != nil {
		return errs.NewRequestError("connection already exists")
	}
	intent := data.ConnectionIntent{
		Type: data.INTENT_TYPE_ASSIGNED,
	}
	createdAt := time.Now()
	mentorship := data.Mentorship{
		MentorUserId: mentorUserId,
		CreatedAt:    createdAt,
	}
	conn := data.Connection{
		UserOneId:  mentorUserId,
		UserTwoId:  menteeUserId,
		CreatedAt:  createdAt,
		AcceptedAt: &createdAt, // Automatically accept.
		Intent:     &intent,
		Mentorship: &mentorship,
	}
	if requestType == api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN {
		rlog.Infof("not a dry run, adding (%d, %d)", mentorUserId, menteeUserId)
		if err := tx.Create(&conn).Error; err != nil {
			return errs.NewDbError(err)
		}
		if err := meetup_reminder.ScheduleInitialReminder(tx, conn.UserOneId, conn.UserTwoId); err != nil {
			return errs.NewDbError(err)
		}
	}
	return nil
}

// More custom method for adding a mentorship used when creating a mentorship using a match round
// If no connection exists, creates connection, mentorship, connection_match_round
// If connection exists, but no mentorship, creates connection
// If connection exists, but no connection_match_round exists for this round, adds
// connection_match_round
func AddMatchRoundMentorship(
	tx *gorm.DB,
	mentorUserId data.TUserID,
	menteeUserId data.TUserID,
	matchRoundId data.TMatchRoundID,
) error {
	var conn data.Connection
	err := tx.Where(
		"user_one_id = ? AND user_two_id = ? OR user_one_id = ? AND user_two_id = ?",
		mentorUserId,
		menteeUserId,
		menteeUserId,
		mentorUserId,
	).Preload("Mentorship").Preload("MatchRounds").Find(&conn).Error

	createdAt := time.Now()
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			conn = data.Connection{
				UserOneId:  mentorUserId,
				UserTwoId:  menteeUserId,
				CreatedAt:  createdAt,
				AcceptedAt: &createdAt, // Automatically accept.
			}
		} else {
			return err
		}
	}

	if conn.Mentorship == nil {
		conn.Mentorship = &data.Mentorship{
			MentorUserId: mentorUserId,
			CreatedAt:    createdAt,
		}
	}

	hasMatchRound := false
	if conn.MatchRounds != nil {
		for _, matchRound := range conn.MatchRounds {
			if matchRound.MatchRoundId == matchRoundId {
				hasMatchRound = true
			}
		}
	}

	err = tx.Save(&conn).Error
	if err != nil {
		return err
	}

	if !hasMatchRound {
		connMatchRound := &data.ConnectionMatchRound{
			ConnectionId: conn.ConnectionId,
			MatchRoundId: matchRoundId,
		}
		err = tx.Save(&connMatchRound).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func SendMentorshipNotifications(db *gorm.DB, mentor, mentee *data.User) errs.Error {
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

	emailErr1 := email.SendNewMenteeEmail(mentorEmail, mentor.FirstName, mentee.FirstName, mentee.Email, menteeCohort.ProgramName, menteeCohort.GradYear)
	emailErr2 := email.SendNewMentorEmail(menteeEmail, mentor.FirstName, mentee.FirstName, mentor.Email, mentorCohort.ProgramName, mentorCohort.GradYear)
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

func sendMentorshipNotificationsHandler(
	db *gorm.DB, request *api.CreateMentorshipByEmail) errs.Error {
	mentor, err := query.GetUserByEmail(db, request.MentorEmail)
	if err != nil {
		return err
	}

	mentee, err := query.GetUserByEmail(db, request.MenteeEmail)
	if err != nil {
		return err
	}
	return SendMentorshipNotifications(db, mentor, mentee)
}
