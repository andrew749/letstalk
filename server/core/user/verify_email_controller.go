package user

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"letstalk/server/email"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"time"
	"regexp"
)

var uwEmailRegex = regexp.MustCompile(".*@(edu.)?uwaterloo.ca$")

func SendEmailVerificationController(c *ctx.Context) errs.Error {
	var req *api.SendAccountVerificationEmailRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	if !uwEmailRegex.MatchString(req.Email) {
		return errs.NewRequestError("Expected valid @edu.uwaterloo.ca or @uwaterloo.ca email address")
	}
	user, _ := query.GetUserById(c.Db, c.SessionData.UserId)

	if err := generateAndSendNewAccountVerificationEmail(c, user, req.Email); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

// Transactionally generates a new VerifyEmailId and sends a link in an email to the user.
func generateAndSendNewAccountVerificationEmail(c *ctx.Context, user *data.User, emailAddr string) errs.Error {
	if user.IsEmailVerified {
		return errs.NewRequestError("Account email has already been verified.")
	}
	dbErr := c.WithinTx(func (tx *gorm.DB) error {
		var verifyEmailId *data.VerifyEmailId
		var err error
		if verifyEmailId, err = generateNewVerifyEmailId(tx, user.UserId); err != nil {
			return err
		}
		if err = sendAccountVerifyEmail(verifyEmailId, user, emailAddr); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	return nil
}

// First parameter should be a db transaction.
func generateNewVerifyEmailId(tx *gorm.DB, userId data.TUserID, emailAddr string) (*data.VerifyEmailId, error) {
	var id = uuid.New()
	verifyEmailData := data.VerifyEmailId{
		Id:             id.String(),
		UserId:         userId,
		Email:          emailAddr,
		IsActive:       true,
		ExpirationDate: time.Now().AddDate(0, 0, 1), // Verification email valid for 24 hours.
	}
	// Set all existing VerifyEmailId entries for this user to inactive.
	err := tx.Where(&data.VerifyEmailId{UserId: userId}).
		Update(data.VerifyEmailId{IsActive: false}).
		Error
	if err != nil {
		return nil, err
	}
	// Insert the new VerifyEmailId entry.
	if err := tx.Save(&verifyEmailData).Error; err != nil {
		return nil, err
	}
	return &verifyEmailData, nil
}

func sendAccountVerifyEmail(requestId *data.VerifyEmailId, user *data.User, emailAddr string) error {
	verifyEmailLink := fmt.Sprintf(
		"%s/verify_email.html?requestid=%s",
		utility.BaseUrl,
		requestId.Id,
	)

	// send email to user with link to verify email address
	to := mail.NewEmail(user.FirstName, emailAddr)
	if err := email.SendAccountVerifyEmail(to, verifyEmailLink); err != nil {
		return err
	}
	return nil
}

// VerifyEmailController verifies a new user account's email.
func VerifyEmailController(c *ctx.Context) errs.Error {
	var req api.VerifyEmailRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	verifyEmailId := data.VerifyEmailId{Id: req.Id}
	if err := c.Db.First(&verifyEmailId).Error; err != nil {
		return errs.NewRequestError("Invalid email verification id")
	}

	user, err := query.GetUserById(c.Db, verifyEmailId.UserId)
	if err != nil {
		return errs.NewInternalError("Invalid user id")
	}

	if user.IsEmailVerified {
		// User email already verified, do nothing.
		return nil
	}

	if !verifyEmailId.IsActive || verifyEmailId.ExpirationDate.Before(time.Now()) {
		return errs.NewRequestError("Email verification link is expired")
	}

	c.WithinTx(func(tx *gorm.DB) error {
		// Set all existing VerifyEmailId entries for this user to inactive.
		err := tx.Where(&data.VerifyEmailId{UserId: verifyEmailId.UserId}).
			Update(data.VerifyEmailId{IsActive: false}).
			Error
		if err != nil {
			return err
		}
		// Mark this verify email id as used.
		if err := tx.Model(verifyEmailId).Update(data.VerifyEmailId{IsUsed: true}).Error; err != nil {
			return err
		}
		// Set user's IsEmailVerified to true.
		if err := tx.Model(&user).Update(data.User{IsEmailVerified: true}).Error; err != nil {
			return err
		}
		return nil
	})

	c.Result = "Ok"
	return nil
}
