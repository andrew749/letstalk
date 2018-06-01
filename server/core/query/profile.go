package query

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func UpdateProfile(db *gorm.DB, userId int, request api.ProfileEditRequest) errs.Error {
	bday := time.Unix(request.Birthdate, 0)

	tx := db.Begin()
	err := tx.Model(&data.User{}).Where(&data.User{
		UserId: userId,
	}).Update(data.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Gender:    request.Gender,
		Birthdate: &bday,
	}).Error
	if err != nil {
		tx.Rollback()
		return errs.NewDbError(err)
	}

	err = tx.Model(&data.UserCohort{}).Where(&data.UserCohort{
		UserId: userId,
	}).Update(data.UserCohort{CohortId: request.CohortId}).Error
	if err != nil {
		tx.Rollback()
		return errs.NewDbError(err)
	}

	if request.PhoneNumber != nil {
		err = tx.Model(&data.ExternalAuthData{}).Where(&data.ExternalAuthData{
			UserId: userId,
		}).Update(data.ExternalAuthData{
			PhoneNumber: request.PhoneNumber,
		}).FirstOrCreate(&data.ExternalAuthData{}).Error
		if err != nil {
			tx.Rollback()
			return errs.NewDbError(err)
		}
	}

	if request.MentorshipPreference != nil || request.Bio != nil || request.Hometown != nil {
		// Should only replace non-null elements.
		err = tx.Where(
			&data.UserAdditionalData{UserId: userId},
		).Assign(
			&data.UserAdditionalData{
				MentorshipPreference: request.MentorshipPreference,
				Bio:                  request.Bio,
				Hometown:             request.Hometown,
			},
		).FirstOrCreate(&data.UserAdditionalData{}).Error
		if err != nil {
			tx.Rollback()
			return errs.NewInternalError(err.Error())
		}
	}

	tx.Commit()
	return nil
}

func GetProfile(db *gorm.DB, userId int) (*api.ProfileResponse, errs.Error) {
	user, err := GetUserProfileById(db, userId)
	if err != nil {
		return nil, errs.NewClientError("Unable to get user data.")
	}
	userCohort, err := GetUserCohort(db, userId)
	if err != nil {
		// TODO: Should probably check what the errors here are. Right now assume that cohort does not
		// exist
	}

	userModel := api.ProfileResponse{
		UserPersonalInfo: api.UserPersonalInfo{
			UserId:    userId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthdate: user.Birthdate.Unix(),
			Secret:    user.Secret,
		},
		UserContactInfo: api.UserContactInfo{
			Email: user.Email,
		},
		UserAdditionalData: api.UserAdditionalData{
			MentorshipPreference: user.AdditionalData.MentorshipPreference,
			Bio:                  user.AdditionalData.Bio,
			Hometown:             user.AdditionalData.Hometown,
		},
	}

	if user.ExternalAuthData != nil {
		userModel.UserContactInfo.PhoneNumber = user.ExternalAuthData.PhoneNumber
		userModel.UserContactInfo.FbId = user.ExternalAuthData.FbUserId
	}

	if userCohort != nil {
		userModel.Cohort.CohortId = userCohort.CohortId
		userModel.Cohort.ProgramId = userCohort.ProgramId
		userModel.Cohort.GradYear = userCohort.GradYear
		userModel.Cohort.SequenceId = userCohort.SequenceId
	}

	return &userModel, nil
}

func GetMatchProfile(
	db *gorm.DB,
	meUserId int,
	matchUserId int,
) (*api.ProfileResponse, errs.Error) {
	return GetProfile(db, matchUserId)
}
