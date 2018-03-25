package push

import (
	"fmt"
	"letstalk/server/core/secrets"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendTestEmail() {
	from := mail.NewEmail("Andrew Codispoti", "andrewcod749@gmail.com")
	to := mail.NewEmail("ANDREW", "andrewcod749@gmail.com")
	subject := "SUBJECT"
	plainTextContent := "BOIIIIII"
	htmlContent := "<strong>BOIIIIII</strong>"
	err := SendEmail(from, to, subject, plainTextContent, htmlContent)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func SendEmail(
	from *mail.Email,
	to *mail.Email,
	subject string,
	plainTextContent string,
	htmlContent string,
) error {
	secretsPath := os.Getenv("SECRETS_PATH")
	secrets.LoadSecrets(secretsPath)
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
