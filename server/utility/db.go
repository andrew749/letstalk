package utility

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// GetDB Gets a connection to the gorm db instance using command line params
func GetDB() (*gorm.DB, error) {
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
