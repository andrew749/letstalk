package test

import (
	"fmt"
	"letstalk/server/data"
	"os"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/romana/rlog"
)

type Test struct {
	Test     func(db *gorm.DB)
	TestName string
}

// DB flags
var (
	databasePrefix = uuid.New().String()
	dbPath         = fmt.Sprintf("/tmp/%s.db", databasePrefix)
)

func createFileIfNotExists(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func GetSqliteDB() (*gorm.DB, error) {
	if err := createFileIfNotExists(dbPath); err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	provisionDatabase(db)

	return db, err
}

func provisionDatabase(db *gorm.DB) {
	data.CreateDB(db)
}

func TearDownLocalDatabase() {
	os.Remove(dbPath)
}

// RunTestsWithDb: Run the following tests and fail if any fail.
func RunTestsWithDb(tests []Test) {
	var db *gorm.DB
	var err error
	TearDownLocalDatabase()
	if db, err = GetSqliteDB(); err != nil {
		rlog.Errorf("Failed to provision db %s", err.Error())
		panic(err)
	}
	rlog.Info("Provisioned DB")
	defer db.Close()

	for _, test := range tests {
		runTestWithDb(db, test)
	}

	TearDownLocalDatabase()
}

func RunTestWithDb(test Test) {
	tests := []Test{test}
	RunTestsWithDb(tests)
}

func runTestWithDb(db *gorm.DB, test Test) {
	test.Test(db)
}
