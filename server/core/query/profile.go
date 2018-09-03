package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
)

// TODO: In the future, probably want to remove updating cohort from here since it's not technically
// part of the core profile.
func UpdateProfile(
	db *gorm.DB,
	es *elastic.Client,
	userId data.TUserID,
	request api.ProfileEditRequest,
) errs.Error {
	var userCohort *data.UserCohort
	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		err := db.Model(&data.User{}).Where(&data.User{
			UserId: userId,
		}).Update(data.User{
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Gender:    request.Gender,
			Birthdate: request.Birthdate,
		}).Error
		if err != nil {
			return err
		}

		userCohort, err = updateUserCohort(db, userId, request.CohortId)
		if err != nil {
			return err
		}

		if request.PhoneNumber != nil {
			err = db.Model(&data.ExternalAuthData{}).Where(&data.ExternalAuthData{
				UserId: userId,
			}).Update(data.ExternalAuthData{
				PhoneNumber: request.PhoneNumber,
			}).FirstOrCreate(&data.ExternalAuthData{}).Error
			if err != nil {
				return err
			}
		}

		if request.MentorshipPreference != nil || request.Bio != nil || request.Hometown != nil {
			// Should only replace non-null elements.
			err = db.Where(
				&data.UserAdditionalData{UserId: userId},
			).Assign(
				&data.UserAdditionalData{
					MentorshipPreference: request.MentorshipPreference,
					Bio:                  request.Bio,
					Hometown:             request.Hometown,
				},
			).FirstOrCreate(&data.UserAdditionalData{}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errs.NewDbError(err)
	}

	go indexCohortMultiTrait(es, userCohort)

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
	}

	if user.AdditionalData != nil {
		userModel.UserAdditionalData.MentorshipPreference = user.AdditionalData.MentorshipPreference
		userModel.UserAdditionalData.Bio = user.AdditionalData.Bio
		userModel.UserAdditionalData.Hometown = user.AdditionalData.Hometown
	}

	if user.ExternalAuthData != nil {
		userModel.UserContactInfo.PhoneNumber = user.ExternalAuthData.PhoneNumber
		userModel.UserContactInfo.FbId = user.ExternalAuthData.FbUserId
		userModel.UserContactInfo.FbLink = user.ExternalAuthData.FbProfileLink
	}

	if user.UserPositions != nil {
		userModel.UserPositions = make([]api.UserPosition, len(user.UserPositions))
		for i, userPosition := range user.UserPositions {
			userModel.UserPositions[i] = api.UserPosition{
				Id:               userPosition.Id,
				RoleId:           userPosition.RoleId,
				RoleName:         userPosition.RoleName,
				OrganizationId:   userPosition.OrganizationId,
				OrganizationName: userPosition.OrganizationName,
				OrganizationType: userPosition.OrganizationType,
				StartDate:        userPosition.StartDate,
				EndDate:          userPosition.EndDate,
			}
		}
	}

	if user.UserSimpleTraits != nil {
		userModel.UserSimpleTraits = make([]api.UserSimpleTrait, len(user.UserSimpleTraits))
		for i, userSimpleTrait := range user.UserSimpleTraits {
			userModel.UserSimpleTraits[i] = api.UserSimpleTrait{
				Id:                     userSimpleTrait.Id,
				SimpleTraitId:          userSimpleTrait.SimpleTraitId,
				SimpleTraitName:        userSimpleTrait.SimpleTraitName,
				SimpleTraitType:        userSimpleTrait.SimpleTraitType,
				SimpleTraitIsSensitive: userSimpleTrait.SimpleTraitIsSensitive,
			}
		}
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
