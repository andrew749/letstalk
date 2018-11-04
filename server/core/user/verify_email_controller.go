package user

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/core/utility/uw_email"
	"letstalk/server/data"
	"letstalk/server/email"

	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailVerificationController(c *ctx.Context) errs.Error {
	var req *api.SendAccountVerificationEmailRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	if err := handleSendAccountVerificationEmailRequest(c, req); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

// Transactionally generates a new VerifyEmailId and sends a link in an email to the user.
func handleSendAccountVerificationEmailRequest(c *ctx.Context, req *api.SendAccountVerificationEmailRequest) errs.Error {
	if !uw_email.Validate(req.Email) {
		return errs.NewRequestError("Expected valid @edu.uwaterloo.ca or @uwaterloo.ca email address")
	}
	uwEmail := uw_email.OfString(req.Email)
	user, _ := query.GetUserById(c.Db, c.SessionData.UserId)

	if user.IsEmailVerified {
		return errs.NewRequestError("Account email has already been verified.")
	}
	// TODO ensure that this email hasn't already been used for a different account
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		var verifyEmailId *data.VerifyEmailId
		var err error
		if verifyEmailId, err = query.GenerateNewVerifyEmailId(tx, user.UserId, uwEmail); err != nil {
			return err
		}
		if err = sendAccountVerifyEmail(verifyEmailId, user, uwEmail); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewInternalError("error sending account verify email: %v", dbErr)
	}
	return nil
}

func sendAccountVerifyEmail(requestId *data.VerifyEmailId, user *data.User, uwEmail uw_email.UwEmail) error {
	verifyEmailLink := fmt.Sprintf(
		"%s/verify_email.html?requestId=%s",
		utility.BaseUrl,
		requestId.Id,
	)

	// send email to user with link to verify email address
	to := mail.NewEmail(user.FirstName, uwEmail.ToStringNormalized())
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
	if err := handleEmailVerification(c, &req); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func handleEmailVerification(c *ctx.Context, req *api.VerifyEmailRequest) errs.Error {
	verifyEmailId := data.VerifyEmailId{Id: req.Id}
	if err := c.Db.First(&verifyEmailId).Error; err != nil {
		return errs.NewRequestError("Invalid email verification id")
	}

	user, err := query.GetUserById(c.Db, verifyEmailId.UserId)
	if err != nil {
		return errs.NewRequestError("Invalid user id")
	}

	if user.IsEmailVerified {
		// User email already verified, do nothing.
		return nil
	}

	if !verifyEmailId.IsActive || verifyEmailId.ExpirationDate.Before(time.Now()) {
		return errs.NewRequestError("Email verification link is expired")
	}

	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		// Set all existing VerifyEmailId entries for this user to inactive.
		err := tx.Model(&data.VerifyEmailId{}).
			Where(&data.VerifyEmailId{UserId: verifyEmailId.UserId}).
			Update("is_active", false).Error
		if err != nil {
			return err
		}
		// Mark this verify email id as used.
		if err := tx.Model(&verifyEmailId).Update("is_used", true).Error; err != nil {
			return err
		}
		// Set user's IsEmailVerified to true.
		if err := tx.Model(&user).Update(data.User{IsEmailVerified: true}).Error; err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	return nil
}
