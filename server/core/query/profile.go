package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func UpdateProfile(db *gorm.DB, userId data.TUserID, request api.ProfileEditRequest) errs.Error {
	tx := db.Begin()
	err := tx.Model(&data.User{}).Where(&data.User{
		UserId: userId,
	}).Update(data.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Gender:    request.Gender,
		Birthdate: request.Birthdate,
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

func GetProfile(db *gorm.DB, userId data.TUserID) (*api.ProfileResponse, errs.Error) {
	user, err := GetUserProfileById(db, userId)
	if err != nil {
		return nil, errs.NewRequestError("Unable to get user data.")
	}
	userCohort, err := GetUserCohort(db, userId)
	if err != nil {
		// TODO: Should probably check what the errors here are. Right now assume that cohort does not
		// exist
	}

	userModel := api.ProfileResponse{
		UserPersonalInfo: api.UserPersonalInfo{
			UserId:     userId,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Gender:     user.Gender,
			Birthdate:  user.Birthdate,
			Secret:     user.Secret,
			ProfilePic: user.ProfilePic,
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
		userModel.UserContactInfo.FbLink = user.ExternalAuthData.FbProfileLink
	}

	if userCohort != nil {
		// NOTE: New API will allow for null sequence ids.
		sequenceId := ""
		if userCohort.SequenceId != nil {
			sequenceId = *userCohort.SequenceId
		}

		userModel.Cohort.CohortId = userCohort.CohortId
		userModel.Cohort.ProgramId = userCohort.ProgramId
		userModel.Cohort.GradYear = userCohort.GradYear
		userModel.Cohort.SequenceId = sequenceId
	}

	return &userModel, nil
}

func GetMatchProfile(
	db *gorm.DB,
	meUserId data.TUserID,
	matchUserId data.TUserID,
) (*api.ProfileResponse, errs.Error) {

	// Fetch mentors and mentees.
	flag := api.MATCHING_INFO_FLAG_NONE
	// Matchings where user is the mentee.
	mentors, err := GetMentorsByMenteeId(db, meUserId, flag)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	// Matchings where user is the mentor.
	mentees, err := GetMenteesByMentorId(db, meUserId, flag)
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	reqFlag := api.REQ_MATCHING_INFO_FLAG_NONE
	// Request matchings where user is answerer.
	askers, err := GetAskersByAnswererId(db, meUserId, reqFlag)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	// Request matchings where user is asker.
	answerers, err := GetAnswerersByAskerId(db, meUserId, reqFlag)
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	userIds := make(map[data.TUserID]interface{})
	for _, mentor := range mentors {
		userIds[mentor.MentorUser.UserId] = nil
	}
	for _, mentee := range mentees {
		userIds[mentee.MenteeUser.UserId] = nil
	}
	for _, asker := range askers {
		userIds[asker.AskerUser.UserId] = nil
	}
	for _, answerer := range answerers {
		userIds[answerer.AnswererUser.UserId] = nil
	}

	// Check if the user profile being request is actually matched with the calling user
	if _, ok := userIds[matchUserId]; !ok {
		return nil, errs.NewRequestError("You are not matched with this user")
	}

	return GetProfile(db, matchUserId)
}
