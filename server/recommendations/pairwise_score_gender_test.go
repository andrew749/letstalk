package recommendations

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
)

func TestGenderPairwiseScoreCalculateMF(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.Gender = data.GENDER_MALE
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			user2.Gender = data.GENDER_FEMALE
			err = db.Save(user2).Error
			assert.NoError(t, err)

			score, err := GenderPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(0.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGenderPairwiseScoreCalculateMM(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.Gender = data.GENDER_MALE
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			user2.Gender = data.GENDER_MALE
			err = db.Save(user2).Error
			assert.NoError(t, err)

			score, err := GenderPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(1.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGenderPairwiseScoreCalculateFF(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.Gender = data.GENDER_FEMALE
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			user2.Gender = data.GENDER_FEMALE
			err = db.Save(user2).Error
			assert.NoError(t, err)

			score, err := GenderPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(1.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGenderPairwiseScoreCalculateUnspecifiedFirst(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.Gender = data.GENDER_UNSPECIFIED
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			user2.Gender = data.GENDER_FEMALE
			err = db.Save(user2).Error
			assert.NoError(t, err)

			score, err := GenderPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(1.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGenderPairwiseScoreCalculateUnspecifiedSecond(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.Gender = data.GENDER_FEMALE
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			user2.Gender = data.GENDER_UNSPECIFIED
			err = db.Save(user2).Error
			assert.NoError(t, err)

			score, err := GenderPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(1.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}
