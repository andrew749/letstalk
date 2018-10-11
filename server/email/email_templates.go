package email

import "github.com/sendgrid/sendgrid-go/helpers/mail"

const (
	SubscribeEmail      = "eaf48eac-ef8a-4dfc-9b10-e09f2dc4b337"
	PasswordChangeEmail = "5bd55885-2793-4848-8ed1-18d2483188f8"
	NewAccount          = "3df07433-f8a3-453f-a94c-1408e85d35e4"
	AccountVerifyEmail  = "0f4be460-b8b2-42bb-8682-222af0ddba99"
	NewMentorEmail      = "d-7d402b5dbdee4bb9b2b94f4eb6e1bdb5"
	NewMenteeEmail      = "d-f082bc47341e40b3ad40c71d2f93621d"
)

func SendSubscribeEmail(
	to *mail.Email,
	name string,
) error {
	message := CreateBasicTemplatedEmail(to, SubscribeEmail, nil)
	return SendEmail(message)
}

func SendForgotPasswordEmail(
	to *mail.Email,
	passwordChangeLink string,
) error {
	var emailContext interface{} = struct {
		RecipientEmail     string `email_sub:":recipientemail"`
		PasswordChangeLink string `email_sub:":passwordchangelink"`
	}{
		to.Address,
		passwordChangeLink,
	}

	message := CreateBasicTemplatedEmail(to, PasswordChangeEmail, &emailContext)

	return SendEmail(message)
}

func SendNewAccountEmail(
	to *mail.Email,
	name string,
) error {
	var emailContext interface{} = struct {
		Name string `email_sub:":name"`
	}{
		Name: name,
	}

	message := CreateBasicTemplatedEmail(to, NewAccount, &emailContext)
	return SendEmail(message)
}

func SendAccountVerifyEmail(
	to *mail.Email,
	verifyEmailLink string,
) error {
	var emailContext interface{} = struct {
		RecipientEmail  string `email_sub:":recipientemail"`
		VerifyEmailLink string `email_sub:":verifyemaillink"`
	}{
		to.Address,
		verifyEmailLink,
	}

	message := CreateBasicTemplatedEmail(to, AccountVerifyEmail, &emailContext)

	return SendEmail(message)
}

func SendNewMentorEmail(
	to *mail.Email,
	mentorName string,
	menteeName string,
	mentorCohort string,
	menteeCohort string,
) error {
	var emailContext interface{} = struct {
		MentorName   string `email_sub:":mentor_name"`
		MenteeName   string `email_sub:":mentee_name"`
		MentorCohort string `email_sub:":mentor_cohort"`
		MenteeCohort string `email_sub:":mentee_cohort"`
	}{}

	message := CreateBasicTemplatedEmail(to, NewMentorEmail, &emailContext)

	return SendEmail(message)
}

func SendNewMenteeEmail(
	to *mail.Email,
	mentorName string,
	menteeName string,
	mentorCohort string,
	menteeCohort string,
) error {
	var emailContext interface{} = struct {
		MentorName   string `email_sub:":mentor_name"`
		MenteeName   string `email_sub:":mentee_name"`
		MentorCohort string `email_sub:":mentor_cohort"`
		MenteeCohort string `email_sub:":mentee_cohort"`
	}{}

	message := CreateBasicTemplatedEmail(to, NewMenteeEmail, &emailContext)

	return SendEmail(message)
}
