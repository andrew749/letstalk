package test_light

// A lightweight package that has no dependencies on components of hive

import (
	"fmt"
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

func createFileIfNotExists(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func GetSqliteDB(dbPath string) (*gorm.DB, error) {
	if err := createFileIfNotExists(dbPath); err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, err
}

func TearDownLocalDatabase(dbPath string) {
	os.Remove(dbPath)
}

type DatabaseAwareFunc func(*gorm.DB) error

// RunTestsWithDb: Run the following tests and fail if any fail.
func RunTestsWithDb(provisionDatabase DatabaseAwareFunc, tests []Test) {
	var db *gorm.DB
	var err error
	databasePrefix := uuid.New().String()

	dbPath := fmt.Sprintf("/tmp/%s.db", databasePrefix)
	TearDownLocalDatabase(dbPath)
	if db, err = GetSqliteDB(dbPath); err != nil {
		rlog.Errorf("Failed to create db %s", err.Error())
		panic(err)
	}
	defer db.Close()

	if err := provisionDatabase(db); err != nil {
		rlog.Errorf("Unable to provision db %+v", err)
	}

	rlog.Info("Provisioned DB")

	for _, test := range tests {
		runTestWithDb(db, test)
	}

	TearDownLocalDatabase(dbPath)
}

func RunTestWithDb(databaseProvision DatabaseAwareFunc, test Test) {
	tests := []Test{test}
	RunTestsWithDb(databaseProvision, tests)
}

func runTestWithDb(db *gorm.DB, test Test) {
	test.Test(db)
}
