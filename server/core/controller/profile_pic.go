package controller

import (
	"bytes"
	"image"
	"image/jpeg"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/s3_assets"
	"letstalk/server/data"
)

// UploadProfilePic Controller to handle api calls to upload profile pictures to s3
func UploadProfilePic(ctx *ctx.Context) errs.Error {
	file, _, err := ctx.GinContext.Request.FormFile("photo")
	if err != nil {
		return errs.NewInternalError("Unable to decode message")
	}

	var imageData image.Image
	if imageData, _, err = image.Decode(file); err != nil {
		return errs.NewInternalError("Unable to decode image")
	}

	buf := new(bytes.Buffer)
	if err = jpeg.Encode(buf, imageData, nil); err != nil {
		return errs.NewInternalError("Unable to encode image")
	}
	reader := bytes.NewReader(buf.Bytes())

	var profilePicLocation *string
	if profilePicLocation, err = s3_assets.UploadProfilePic(
		ctx.SessionData.UserId,
		reader,
	); err != nil {
		return errs.NewInternalError("Unable to upload image")
	}

	// TODO(wojtechnology): Move into deeper layer so we don't break abstraction boundary
	db := ctx.Db
	var user *data.User
	// update user profile pic locatio in data
	if user, err = query.GetUserById(db, ctx.SessionData.UserId); err != nil {
		return errs.NewInternalError("Unable to update image")
	}

	if err := db.Model(*user).Updates(data.User{ProfilePic: profilePicLocation}).Error; err != nil {
		return errs.NewInternalError("Unable to update image")
	}

	ctx.Result = &api.ProfilePicUploadResult{Status: "ok", NewLocation: profilePicLocation}
	return nil
}
