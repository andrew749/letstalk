package push

import (
	"fmt"
	"letstalk/server/core/secrets"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(
	from *mail.Email,
	to *mail.Email,
	subject string,
	plainTextContent string,
	htmlContent string,
) error {
	message := mail.NewSingleEmail(
		from,
		subject,
		to,
		plainTextContent,
		htmlContent,
	)
	fmt.Println(secrets.GetSecrets().SendGrid)
	client := sendgrid.NewSendClient(secrets.GetSecrets().SendGrid)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

func SendForgotPasswordEmail(
	to *mail.Email,
	passwordChangeLink string,
) error {
	client := sendgrid.NewSendClient(secrets.GetSecrets().SendGrid)
	// Create message and configure with empty text to force template
	// body to be used instead. You have to set up a transactional
	// template on SendGrid's web site and reference its ID below where
	// it says <template_id>.
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail("Hive", "andrew@hiveapp.org"))

	// personalize the email to the user
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	// set recipients
	p.AddTos(tos...)
	p.SetSubstitution(":passwordchangelink", passwordChangeLink)

	message.AddPersonalizations(p)
	message.SetTemplateID(PasswordChangeEmail)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return err

}

func SendSubscribeEmail(
	to *mail.Email,
	name string,
) error {
	client := sendgrid.NewSendClient(secrets.GetSecrets().SendGrid)
	// Create message and configure with empty text to force template
	// body to be used instead. You have to set up a transactional
	// template on SendGrid's web site and reference its ID below where
	// it says <template_id>.
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail("Andrew@Hive", "andrew@hiveapp.org"))

	// personalize the email to the user
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	// set recipients
	p.AddTos(tos...)
	p.SetSubstitution("%Name%", name)

	message.AddPersonalizations(p)
	message.SetTemplateID(SubscribeEmail)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return err
}
