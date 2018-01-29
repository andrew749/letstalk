package users

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/db"
	"letstalk/server/data"
	"log"

	"letstalk/server/core/errs"

	"github.com/mijia/modelq/gmq"
)

func genderIdByName(gender string) int {
	switch gender {
	case "FEMALE":
		return 1
	case "MALE":
		return 2
	default:
		return 0
	}
}

func PostUser(c *ctx.Context) errs.Error {
	var inputUser api.User
	err := c.GinContext.BindJSON(&inputUser)
	if err != nil {
		return errs.NewClientError("failed to bind user: %s", err)
	}
	log.Print("post user: ", inputUser)
	// Check that no user exists with this email.
	existingUser, err := data.UserObjs.Select().Where(data.UserObjs.FilterEmail("=", inputUser.Email)).List(c.Db)
	if err != nil {
		return errs.NewDbError(err)
	}
	if len(existingUser) != 0 {
		return errs.NewClientError("a user already exists with email: %s", inputUser.Email)
	}
	// Look up the existing cohort.
	cohorts, err := data.CohortObjs.Select().Where(
		data.CohortObjs.FilterSequence("=", inputUser.Sequence).
			And(data.CohortObjs.FilterGradYear("=", inputUser.GraduatingYear)).
			And(data.CohortObjs.FilterProgramId("=", inputUser.Program))).
		List(c.Db)
	if err != nil {
		return errs.NewDbError(err)
	}
	if len(cohorts) == 0 {
		return errs.NewClientError("cohort not found")
	}
	// Create user and cohort data structures.

	user := data.User{
		Email:     inputUser.Email,
		Nickname:  inputUser.Nickname,
		Name:      inputUser.FullName,
		Gender:    genderIdByName(inputUser.Gender),
		Birthdate: inputUser.Birthday,
	}
	if user.UserId, err = db.NumId(c); err != nil {
		return errs.NewDbError(err)
	}
	userCohort := data.UserCohort{
		UserId:   user.UserId,
		CohortId: cohorts[0].CohortId,
	}
	// Insert data structures within a transaction.
	dbErr := gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
		if _, err := user.Insert(tx); err != nil {
			return err
		}
		if _, err := userCohort.Insert(tx); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	c.Result = inputUser
	return nil
}
