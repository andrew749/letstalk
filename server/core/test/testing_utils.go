package test

import (
	"letstalk/server/core/test_light"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Heavyweight version of testing utilities that provisions the hive environment

type Test test_light.Test

func provisionDatabase(db *gorm.DB) error {
	data.CreateDB(db)
	// TODO(acod): refactor database spinup to return error instead of panicing
	return nil
}

func convertHeavyTestToLightTest(test Test) test_light.Test {
	return test_light.Test(test)
}

func convertHeavyTestSliceToLightTestSlice(tests []Test) []test_light.Test {
	res := make([]test_light.Test, 0)
	for _, test := range tests {
		res = append(res, convertHeavyTestToLightTest(test))
	}
	return res
}

// RunTestsWithDb: Run the following tests and fail if any fail.
func RunTestsWithDb(tests []Test) {
	test_light.RunTestsWithDb(provisionDatabase, convertHeavyTestSliceToLightTestSlice(tests))
}

func RunTestWithDb(test Test) {
	tests := []Test{test}
	RunTestsWithDb(tests)
}
