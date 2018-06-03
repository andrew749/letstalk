package utility

import (
	"database/sql"
	"flag"
	"fmt"
	"letstalk/server/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Test struct {
	Test     func(db *gorm.DB)
	TestName string
}

// DB flags
var (
	dbUser   = flag.String("db_user", "letstalk", "mySQL user")
	dbPass   = flag.String("db_pass", "uwletstalk", "mySQL password")
	rootUser = "root" // obviously this is just debug
)

// var databasePrefix = uuid.New().String()
var databasePrefix = "integration_test"
var databaseRootConnectionString = fmt.Sprintf("%s:%s@tcp(:3306)/", rootUser, *dbPass)
var databaseConnectionString = fmt.Sprintf("%s:%s@tcp(:3306)/", *dbUser, *dbPass)

func CreateLocalDatabase() (*gorm.DB, error) {
	var err error
	var dbInit *sql.DB
	if dbInit, err = sql.Open("mysql", databaseRootConnectionString); err != nil {
		panic(err)
	}
	// create temp database
	if _, err = dbInit.Exec("CREATE DATABASE " + databasePrefix); err != nil {
		panic(err)
	}

	if _, err = dbInit.Exec("USE " + databasePrefix); err != nil {
		panic(err)
	}

	if _, err = dbInit.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%'", databasePrefix, *dbUser)); err != nil {
		panic(err)
	}

	if _, err = dbInit.Exec("FLUSH PRIVILEGES"); err != nil {
		panic(err)
	}

	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s%s?parseTime=true", databaseConnectionString, databasePrefix),
	)

	if err != nil {
		return nil, err
	}

	data.CreateDB(db)
	return db, nil
}

func TearDownLocalDatabase() {
	var dbInit *sql.DB
	var err error

	if dbInit, err = sql.Open("mysql", databaseRootConnectionString); err != nil {
		panic(err)
	}
	if _, err = dbInit.Exec("DROP DATABASE IF EXISTS " + databasePrefix); err != nil {
		panic(err)
	}
}

// RunTestsWithDb: Run the following tests and fail if any fail.
func RunTestsWithDb(tests []Test) {
	var db *gorm.DB
	var err error
	TearDownLocalDatabase()
	if db, err = CreateLocalDatabase(); err != nil {
		panic(err)
	}

	for _, test := range tests {
		runTestWithDb(db, test)
	}

	TearDownLocalDatabase()
}

func runTestWithDb(db *gorm.DB, test Test) {
	test.Test(db)
}
