package onboarding

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/romana/rlog"
	"letstalk/server/aws_utils"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type ProfilePicUploadRequest struct {
}

type ProfilePicUploadResult struct {
	Status string
}

const (
	PROFILE_PIC_BUCKET = "hive-user-profile-pictures"
)

func ProfilePicController(ctx *ctx.Context) errs.Error {
	data, err := ctx.GinContext.GetRawData()
	if err != nil {
		return errs.NewInternalError("Unable to decode message")
	}

	s3_client, err := aws_utils.GetS3ServiceClient()
	if err != nil {
		return errs.NewInternalError("Unable to connect to s3")
	}
	pictureId := fmt.Sprintf("%d", ctx.SessionData.UserId)
	rlog.Debug("Uploading profile picture file for %d", pictureId)
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(s3_client)

	reader := bytes.NewReader(data)

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(PROFILE_PIC_BUCKET),
		Key:    aws.String(pictureId), // user id as key
		Body:   reader,
	})

	if err != nil {
		return errs.NewInternalError(fmt.Sprintf("failed to upload file, %v", err))
	}

	rlog.Debug("Successfully uploaded profile picture under %d", pictureId)
	ctx.Result = &ProfilePicUploadResult{"ok"}
	return nil
}
