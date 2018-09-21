package query

import (
	"fmt"

	"letstalk/server/core/api"
	"letstalk/server/core/converters"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

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
		res.Cohort = converters.ApiCohortV2FromDataCohort(user.Cohort.Cohort)
	}
	return res
}

func buildUserSearchResponse(isAnonymous bool, users []data.User) *api.UserSearchResponse {
	var results []api.UserSearchResult
	if isAnonymous {
		results = []api.UserSearchResult{}
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

func searchUsersCommon(query *gorm.DB, size int, extraCols *string) *gorm.DB {
	cols := "user_id"
	if extraCols != nil {
		cols = fmt.Sprintf("%s, %s", cols, *extraCols)
	}

	return query.Select(fmt.Sprintf("DISTINCT %s", cols)).Preload(
		"User.Cohort.Cohort",
	).Limit(size).Order("RAND()")
}

func SearchUsersByCohort(
	db *gorm.DB,
	req api.CohortUserSearchRequest,
	userId data.TUserID,
) (*api.UserSearchResponse, errs.Error) {
	var userCohorts []data.UserCohort

	query := db.Where(&data.UserCohort{CohortId: req.CohortId}).Not(&data.UserCohort{
		UserId: userId,
	})
	if err := searchUsersCommon(query, req.Size, nil).Find(&userCohorts).Error; err != nil {
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
	userId data.TUserID,
) (*api.UserSearchResponse, errs.Error) {

	var userSimpleTraits []data.UserSimpleTrait

	extraCols := "simple_trait_is_sensitive"
	query := db.Where(&data.UserSimpleTrait{
		SimpleTraitId: req.SimpleTraitId,
	}).Not(&data.UserSimpleTrait{UserId: userId})
	if err := searchUsersCommon(
		query,
		req.Size,
		&extraCols,
	).Find(&userSimpleTraits).Error; err != nil {
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
	userId data.TUserID,
) (*api.UserSearchResponse, errs.Error) {
	var userPositions []data.UserPosition

	query := db.Where(&data.UserPosition{
		RoleId:         req.RoleId,
		OrganizationId: req.OrganizationId,
	}).Not(&data.UserPosition{UserId: userId})
	if err := searchUsersCommon(query, req.Size, nil).Find(&userPositions).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	users := make([]data.User, len(userPositions))
	for i, userPosition := range userPositions {
		users[i] = *userPosition.User
	}

	return buildUserSearchResponse(false, users), nil
}
