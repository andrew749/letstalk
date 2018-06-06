package utility

import (
	"letstalk/server/data"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Test struct {
	Test     func(db *gorm.DB)
	TestName string
}

// DB flags
var (
	dbPath = "/tmp/test.db"
)

// var databasePrefix = uuid.New().String()
var databasePrefix = "integration_test"

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

	return db, nil
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
		panic(err)
	}
	defer db.Close()

	for _, test := range tests {
		runTestWithDb(db, test)
	}

	TearDownLocalDatabase()
}

func runTestWithDb(db *gorm.DB, test Test) {
	test.Test(db)
}
