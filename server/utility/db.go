package utility

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
	"github.com/romana/rlog"

	"letstalk/server/core/errs"
)

// GetDB Gets a connection to the gorm db instance using command line params
func GetDB() (*gorm.DB, error) {
	flag.Parse()
	rlog.Infof("Connecting to db: %s", *dbAddr)
	return gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@%s/letstalk?charset=utf8mb4&parseTime=true", *dbUser, *dbPass, *dbAddr),
	)
}

// RunWithDb Wrap a call and get a db instance to run the callback with
func RunWithDb(c func(tx *gorm.DB) error) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return c(db)
}

// Correctly wraps an error from a gorm.DB call with an errs.Error, returning a NotFoundError
// if the item is not found and a DBError otherwise.
func WrapGormDBError(dbErr error, notFoundMessage string) errs.Error {
	if dbErr != nil {
		if gorm.IsRecordNotFoundError(dbErr) {
			return errs.NewNotFoundError(notFoundMessage)
		} else {
			return errs.NewDbError(dbErr)
		}
	}
	return nil
}
