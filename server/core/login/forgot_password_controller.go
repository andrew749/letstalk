package login

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/auth"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"letstalk/server/email"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func generateNewForgotPasswordRequest(db *gorm.DB, userId int) (*data.ForgotPasswordId, error) {
	var id = uuid.New()
	forgotPasswordRequest := data.ForgotPasswordId{
		Id:     id.String(),
		UserId: userId,
	}
	if err := db.Save(&forgotPasswordRequest).Error; err != nil {
		return nil, err
	}
	return &forgotPasswordRequest, nil
}

func sendForgotPasswordEmail(db *gorm.DB, requestId *data.ForgotPasswordId, user *data.User) error {
	passwordChangeLink := fmt.Sprintf(
		"%s/change_password.html?requestId=%s",
		utility.BaseUrl,
		requestId.Id,
	)

	// send email to user with link to change password
	to := mail.NewEmail(user.FirstName, user.Email)
	if err := email.SendForgotPasswordEmail(to, passwordChangeLink); err != nil {
		return err
	}
	return nil
}

func GenerateNewForgotPasswordRequestController(ctx *ctx.Context) errs.Error {
	db := ctx.Db

	var err error
	var req *api.GenerateForgotPasswordRequest
	if err = ctx.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	var forgotPasswordId *data.ForgotPasswordId

	var user *data.User
	if user, err = query.GetUserByEmail(db, req.Email); err != nil {
		// return errs.NewRequestError("Can not find a user with that email")
		// this user email does not exist
		return nil
	}

	if forgotPasswordId, err = generateNewForgotPasswordRequest(db, user.UserId); err != nil {
		return errs.NewRequestError(err.Error())
	}

	if err := sendForgotPasswordEmail(db, forgotPasswordId, user); err != nil {
		return errs.NewInternalError(err.Error())
	}

	ctx.Result = "Ok"
	return nil
}

// ForgotPasswordController: has the ability to change a user password unauthenticated
func ForgotPasswordController(ctx *ctx.Context) errs.Error {
	var forgotPasswordRequestChangeId api.ForgotPasswordChangeRequest
	if err := ctx.GinContext.BindJSON(&forgotPasswordRequestChangeId); err != nil {
		return errs.NewRequestError(err.Error())
	}
	tx := ctx.Db.Begin()

	// find the user we need to change
	forgotPasswordId := data.ForgotPasswordId{
		Id: forgotPasswordRequestChangeId.ForgotPasswordRequestId,
	}

	if err := tx.First(&forgotPasswordId).Error; err != nil {
		tx.Rollback()
		return errs.NewRequestError("Invalid password change token")
	}

	if forgotPasswordId.Used {
		tx.Rollback()
		return errs.NewRequestError("Password change token already used.")
	}

	if err := auth.ChangeUserPassword(
		tx,
		forgotPasswordId.UserId,
		forgotPasswordRequestChangeId.NewPassword,
	); err != nil {
		tx.Rollback()
		return errs.NewInternalError(err.Error())
	}

	// update the token to be used
	forgotPasswordId.Used = true
	if err := tx.Save(&forgotPasswordId).Error; err != nil {
		return errs.NewDbError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewDbError(err)
	}
	ctx.Result = "Ok"
	return nil
}
