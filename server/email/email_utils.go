package email

import (
	"letstalk/server/core/secrets"
	"reflect"

	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// general user data passed to an email that is available in the sendgrid template
const (
	substitutionTag      = "email_sub"
	defaultSenderName    = "Hive"
	defaultSenderAddress = "uwhive@gmail.com"
)

// example of how to create a struct with appropriate tags
type BaseUserEmailContext struct {
	Name  string `email_sub:"Name"`
	Email string `email_sub:"Email"`
}

// MarshallEmailSubstitutions: create map of substitutions to make
func MarshallEmailSubstitutions(c interface{}) map[string]string {
	t := reflect.TypeOf(c)
	substitutions := make(map[string]string)
	numFields := t.NumField()
	// go over all fields	and marshall in our map
	for i := 0; i < numFields; i++ {
		tempField := t.Field(i)
		tag := tempField.Tag.Get(substitutionTag)
		substitutions[tag] = getField(c, tempField.Name)
	}
	return substitutions
}

func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

// GetBasicTemplatedEmail
func CreateBasicTemplatedEmail(
	recipient *mail.Email,
	templateID string,
	emailContext *interface{},
) *mail.SGMailV3 {
	// Create message and configure with empty text to force template
	// body to be used instead. You have to set up a transactional
	// template on SendGrid's web site and reference its ID below where
	// it says <template_id>.
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail(defaultSenderName, defaultSenderAddress))

	// personalize the email to the user
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		recipient,
	}
	// set recipients
	p.AddTos(tos...)

	var context interface{}
	if emailContext == nil {
		context = BaseUserEmailContext{
			Name:  recipient.Name,
			Email: recipient.Address,
		}
	} else {
		context = emailContext
	}

	p.Substitutions = MarshallEmailSubstitutions(context)

	message.AddPersonalizations(p)
	message.SetTemplateID(templateID)
	return message
}

// SendBasicEmail  Sends a basic html email.
func SendBasicEmail(
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
	client := sendgrid.NewSendClient(secrets.GetSecrets().SendGrid)
	response, err := client.Send(message)
	if err != nil {
		rlog.Error(err)
	} else {
		rlog.Debug(response.StatusCode)
		rlog.Debug(response.Headers)
	}
	return nil
}

// SendEmail  Helper to deliver an already created email message
func SendEmail(message *mail.SGMailV3) error {
	client := sendgrid.NewSendClient(secrets.GetSecrets().SendGrid)
	response, err := client.Send(message)
	if err != nil {
		rlog.Error(err)
	} else {
		rlog.Debug(response.StatusCode)
		rlog.Debug(response.Headers)
	}
	return err
}
