package email

import "github.com/sendgrid/sendgrid-go/helpers/mail"

const (
	SubscribeEmail      = "eaf48eac-ef8a-4dfc-9b10-e09f2dc4b337"
	PasswordChangeEmail = "5bd55885-2793-4848-8ed1-18d2483188f8"
	NewAccount          = "3df07433-f8a3-453f-a94c-1408e85d35e4"
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
