package s3_assets

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/romana/rlog"

	"letstalk/server/aws_utils"
	"letstalk/server/core/errs"
	"letstalk/server/data"
)

const (
	profilePicBucket = "hive-user-profile-pictures"
)

// UploadProfilePic uploads profile pic for userId to s3 and returns url
func UploadProfilePic(userID data.TUserID, dataReader io.Reader) (*string, error) {
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
