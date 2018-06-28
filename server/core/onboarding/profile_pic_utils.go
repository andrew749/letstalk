package onboarding

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"letstalk/server/aws_utils"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/romana/rlog"
)

// ProfilePicUploadResult Result returned from api controller
type ProfilePicUploadResult struct {
	Status      string
	NewLocation *string
}

const (
	profilePicBucket = "hive-user-profile-pictures"
)

// UploadProfilePic uploads profile pic for userId to s3 and returns url
func UploadProfilePic(userID int, dataReader io.Reader) (*string, error) {
	var s3Client *s3.S3
	var err error

	if s3Client, err = aws_utils.GetS3ServiceClient(); err != nil {
		return nil, err
	}

	// get a random uuid for profile pic
	profilePicID := uuid.New()
	pictureID := fmt.Sprintf("%s", profilePicID)
	rlog.Debugf("Uploading profile picture file for %d", pictureID)
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(s3Client)

	// Upload the file to S3.
	var res *s3manager.UploadOutput
	if res, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(profilePicBucket),
		Key:    aws.String(pictureID), // user id as key
		Body:   dataReader,
	}); err != nil {
		return nil, errs.NewInternalError("failed to upload file, %v", err)
	}

	rlog.Debug("Successfully uploaded profile picture under ", pictureID)
	return &res.Location, nil
}

// ProfilePicController Controller to handle api calls to upload profile pictures to s3
func ProfilePicController(ctx *ctx.Context) errs.Error {
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
	if profilePicLocation, err = UploadProfilePic(ctx.SessionData.UserId, reader); err != nil {
		return errs.NewInternalError("Unable to upload image")
	}

	db := ctx.Db
	var user *data.User
	// update user profile pic locatio in data
	if user, err = query.GetUserById(db, ctx.SessionData.UserId); err != nil {
		return errs.NewInternalError("Unable to update image")
	}

	if err := db.Model(*user).Updates(data.User{ProfilePic: profilePicLocation}).Error; err != nil {
		return errs.NewInternalError("Unable to update image")
	}

	ctx.Result = &ProfilePicUploadResult{Status: "ok", NewLocation: profilePicLocation}
	return nil
}
