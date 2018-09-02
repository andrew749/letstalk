package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func cohortV2FromDataCohort(cohort *data.Cohort) *api.CohortV2 {
	return &api.CohortV2{
		CohortId:     cohort.CohortId,
		ProgramId:    cohort.ProgramId,
		ProgramName:  cohort.ProgramName,
		IsCoop:       cohort.IsCoop,
		GradYear:     cohort.GradYear,
		SequenceId:   cohort.SequenceId,
		SequenceName: cohort.SequenceName,
	}
}

func userSearchResultFromDataUser(user *data.User) *api.UserSearchResult {
	// Reasons can be added later.
	res := &api.UserSearchResult{
		UserId:     user.UserId,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Gender:     user.Gender,
		ProfilePic: user.ProfilePic,
	}

	if user.Cohort != nil && user.Cohort.Cohort != nil {
		res.Cohort = cohortV2FromDataCohort(user.Cohort.Cohort)
	}
	return res
}

func buildUserSearchResponse(isAnonymous bool, users []data.User) *api.UserSearchResponse {
	var results []api.UserSearchResult
	if isAnonymous {
		results = make([]api.UserSearchResult, 0)
	} else {
		results = make([]api.UserSearchResult, len(users))
		for i, user := range users {
			results[i] = *userSearchResultFromDataUser(&user)
		}
	}

	return &api.UserSearchResponse{
		IsAnonymous: isAnonymous,
		NumResults:  len(users),
		Results:     results,
	}
}

func searchUsersCommon(query *gorm.DB, size int) *gorm.DB {
	return query.Preload("User.Cohort.Cohort").Limit(size)
}

func SearchUsersByCohort(
	db *gorm.DB,
	req api.CohortUserSearchRequest,
) (*api.UserSearchResponse, errs.Error) {
	var userCohorts []data.UserCohort

	query := db.Where(&data.UserCohort{CohortId: req.CohortId})
	if err := searchUsersCommon(query, req.Size).Find(&userCohorts).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	users := make([]data.User, len(userCohorts))
	for i, userCohort := range userCohorts {
		users[i] = *userCohort.User
	}

	return buildUserSearchResponse(false, users), nil
}

func SearchUsersBySimpleTrait(
	db *gorm.DB,
	req api.SimpleTraitUserSearchRequest,
) (*api.UserSearchResponse, errs.Error) {

	var userSimpleTraits []data.UserSimpleTrait

	query := db.Where(&data.UserSimpleTrait{SimpleTraitId: req.SimpleTraitId})
	if err := searchUsersCommon(query, req.Size).Find(&userSimpleTraits).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	users := make([]data.User, len(userSimpleTraits))
	for i, userSimpleTrait := range userSimpleTraits {
		users[i] = *userSimpleTrait.User
	}

	isAnonymous := false
	if len(userSimpleTraits) > 0 {
		isAnonymous = userSimpleTraits[0].SimpleTraitIsSensitive
	}

	return buildUserSearchResponse(isAnonymous, users), nil
}

func SearchUsersByPosition(
	db *gorm.DB,
	req api.PositionUserSearchRequest,
) (*api.UserSearchResponse, errs.Error) {
	var userPositions []data.UserPosition

	query := db.Where(&data.UserPosition{RoleId: req.RoleId, OrganizationId: req.OrganizationId})
	if err := searchUsersCommon(query, req.Size).Find(&userPositions).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	users := make([]data.User, len(userPositions))
	for i, userPosition := range userPositions {
		users[i] = *userPosition.User
	}

	return buildUserSearchResponse(false, users), nil
}
