package users

import (
	"errors"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/db"
	"letstalk/server/data"
	"log"

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

func PostUser(c *ctx.Context) {
	var inputUser api.User
	err := c.GinContext.BindJSON(&inputUser)
	if err != nil {
		c.GinContext.Error(err)
		return
	}
	log.Print("post user: ", inputUser)
	// Check that no user exists with this email.
	existingUser, err := data.UserObjs.Select().Where(data.UserObjs.FilterEmail("=", inputUser.Email)).List(c.Db)
	if err != nil {
		c.GinContext.Error(err)
		return
	}
	if len(existingUser) != 0 {
		c.GinContext.Error(errors.New("a user already exists with email: " + inputUser.Email))
		return
	}
	// Look up the existing cohort.
	cohorts, err := data.CohortObjs.Select().Where(
		data.CohortObjs.FilterSequence("=", inputUser.Sequence).
			And(data.CohortObjs.FilterGradYear("=", inputUser.GraduatingYear)).
			And(data.CohortObjs.FilterProgramId("=", inputUser.Program))).
		List(c.Db)
	if err != nil {
		c.GinContext.Error(err)
		return
	}
	if len(cohorts) == 0 {
		c.GinContext.Error(errors.New("cohort not found"))
		return
	}
	// Create user and cohort data structures.
	user := data.User{
		UserId:    db.NumId(c),
		Email:     inputUser.Email,
		Nickname:  inputUser.Nickname,
		Name:      inputUser.FullName,
		Gender:    genderIdByName(inputUser.Gender),
		Birthdate: inputUser.Birthday,
	}
	userCohort := data.UserCohort{
		UserId:   user.UserId,
		CohortId: cohorts[0].CohortId,
	}
	// TODO(aklen): nicer checks for errors in context
	if len(c.GinContext.Errors) > 0 {
		return
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
		c.GinContext.Error(dbErr)
		return
	}
	c.Result = inputUser
}
