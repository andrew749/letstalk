package main

import (
	"fmt"
	"letstalk/server/core/secrets"
	"letstalk/server/push"
	"os"

	"github.com/namsral/flag"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendTestEmail() {
	from := mail.NewEmail("Andrew Codispoti", "andrewcod749@gmail.com")
	to := mail.NewEmail("ANDREW", "andrewcod749@gmail.com")
	subject := "SUBJECT"
	plainTextContent := "BOIIIIII"
	htmlContent := "<strong>BOIIIIII</strong>"
	err := push.SendEmail(from, to, subject, plainTextContent, htmlContent)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func SendTestSubscriptionEmail() {
	recipient := mail.NewEmail("Andrew Codispoti", "andrewcod749@gmail.com")
	name := "Andrew"
	err := push.SendSubscribeEmail(recipient, name)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

var (
	testSubscription = flag.Bool("subscribe", false, "Whether to send a subscription email")
	email            = flag.String("email", "", "email to send to")
	name             = flag.String("name", "", "name to use in the email")
)

func main() {
	// preload secrets so we can send using api
	secretsPath := os.Getenv("SECRETS_PATH")
	secrets.LoadSecrets(secretsPath)
	flag.Parse()

	// SendTestEmail()
	if *testSubscription {
		SendTestSubscriptionEmail()
	}
}
