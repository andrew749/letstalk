package bootstrap

import (
	"database/sql"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
	"letstalk/server/data"
)

type BootstrapState string

/**
 * These states will likely change.
 * Current a later state implies that the previous states are satisfied
 * This is currently a linear state hierarchy
 */
const (
	ACCOUNT_CREATED BootstrapState = "account_created" // first state
	ACCOUNT_SETUP   BootstrapState = "account_setup"   // the account has enough information to proceed
	ACCOUNT_MATCHED BootstrapState = "account_matched" // account has been matched a peer
)

type BootstrapUserRelationshipDataModel struct {
	User      int          `json:"userId" binding:"required"`
	UserType  api.UserType `json:"userType" binding:"required"`
	FirstName string       `json:"firstName" binding:"required"`
	LastName  string       `json:"lastName" binding:"required"`
	Email     string       `json:"email" binding:"required"`
}

type BootstrapResponse struct {
	State            BootstrapState                       `json:"state" binding:"required"`
	Relationships    []BootstrapUserRelationshipDataModel `json:"relationships" binding:"required"`
	Cohort           *data.Cohort                         `json:"cohort" binding:"required"`
	Me               *data.User                           `json:"me" binding:"required"`
	OnboardingStatus *onboarding.OnboardingStatus         `json:"onboardingStatus" binding:"required"`
}

func sqlResultToBoostrapUserRelationshipDataModel(
	relationships *[]BootstrapUserRelationshipDataModel,
	res *sql.Rows,
	userType api.UserType,
) error {
	var (
		id        int
		firstName string
		lastName  string
		email     string
	)

	for res.Next() {
		// bind this row to variables
		err := res.Scan(&id, &firstName, &lastName, &email)
		if err != nil {
			return err
		}
		*relationships = append(
			*relationships,
			BootstrapUserRelationshipDataModel{id, userType, firstName, lastName, email},
		)
	}
	return nil
}

/**
 * Returns what the current status of a user is
 */
func GetCurrentUserBoostrapStatusController(c *ctx.Context) errs.Error {
	user, err := api.GetUserWithId(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	// since this method is authenticated the account needs to exist.
	var response = BootstrapResponse{
		State: ACCOUNT_CREATED,
		Me:    user,
	}

	onboardingInfo, err := onboarding.GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	response.Cohort = onboardingInfo.UserCohort
	response.OnboardingStatus = &onboarding.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}

	if onboardingInfo.State != onboarding.ONBOARDING_DONE {
		// Onboarding not done. We don't need to get relationships.
		c.Result = response
		return nil
	} else {
		response.State = ACCOUNT_SETUP
	}

	println(string(onboardingInfo.State))

	relationships := make([]BootstrapUserRelationshipDataModel, 0)

	find_mentees_statement, err := c.Db.Prepare(
		`	SELECT
				mentee, first_name, last_name, email
			FROM
				matchings
			INNER JOIN
				user ON user.user_id=matchings.mentee
			WHERE
				mentor=?`,
	)

	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	res, err := find_mentees_statement.Query(c.SessionData.UserId)
	if err != nil {
		return errs.NewInternalError("Unable to create statement")
	}

	defer res.Close()
	sqlResultToBoostrapUserRelationshipDataModel(&relationships, res, api.USER_TYPE_MENTEE)
	err = res.Err()

	if err != nil {
		return errs.NewDbError(err)
	}

	find_mentors_statement, err := c.Db.Prepare(
		`	SELECT
				mentor, first_name, last_name, email
			FROM
				matchings
			INNER JOIN
				user ON user.user_id=matchings.mentor
			WHERE
				mentee=?`,
	)

	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	res, err = find_mentors_statement.Query(c.SessionData.UserId)
	if err != nil {
		return errs.NewClientError("Unable to create statement")
	}

	defer res.Close()
	sqlResultToBoostrapUserRelationshipDataModel(&relationships, res, api.USER_TYPE_MENTOR)
	err = res.Err()

	if err != nil {
		return errs.NewDbError(err)
	}

	if len(relationships) > 0 {
		response.State = ACCOUNT_MATCHED
	}

	response.Relationships = relationships
	c.Result = response

	return nil
}
