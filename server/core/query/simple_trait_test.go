package query

import (
	"letstalk/server/core/test"
	"letstalk/server/data"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddUserSimpleTraitByName(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			err = AddUserSimpleTraitByName(db, nil, data.TUserID(1), "Cycling ")
			assert.NoError(t, err)

			var traits []data.UserSimpleTrait
			err = db.Where(&data.UserSimpleTrait{UserId: 1}).Preload("SimpleTrait").Find(&traits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(traits))
			assert.Equal(t, "Cycling", traits[0].SimpleTraitName)
			assert.Equal(t, data.SIMPLE_TRAIT_TYPE_UNDETERMINED, traits[0].SimpleTraitType)
			assert.False(t, traits[0].SimpleTraitIsSensitive)
			assert.NotNil(t, traits[0].SimpleTrait)
			assert.Equal(t, "Cycling", traits[0].SimpleTrait.Name)
			assert.True(t, traits[0].SimpleTrait.IsUserGenerated)
		},
		TestName: "Test adding simple trait by name",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitByNameAlreadyExists(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{
				Name:            "Cycling",
				IsUserGenerated: false,
			}
			err := db.Save(&trait).Error
			assert.NoError(t, err)

			var traits []data.SimpleTrait
			err = db.Find(&traits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(traits))

			err = AddUserSimpleTraitByName(db, nil, data.TUserID(1), "Cycling ")
			assert.NoError(t, err)

			db.Find(&traits)
			assert.Equal(t, 1, len(traits))

			var userTraits []data.UserSimpleTrait
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("SimpleTrait").Find(&userTraits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userTraits))
			assert.Equal(t, "Cycling", userTraits[0].SimpleTraitName)
			assert.Equal(t, trait, *userTraits[0].SimpleTrait)
		},
		TestName: "Test adding simple trait by name, trait already exists",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitByNameAlreadyExistsIgnoreCase(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{
				Name:            "Cycling",
				IsUserGenerated: false,
			}
			err := db.Save(&trait).Error
			assert.NoError(t, err)

			var traits []data.SimpleTrait
			err = db.Find(&traits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(traits))

			err = AddUserSimpleTraitByName(db, nil, data.TUserID(1), "cycling ")
			assert.NoError(t, err)

			db.Find(&traits)
			assert.Equal(t, 1, len(traits))

			var userTraits []data.UserSimpleTrait
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("SimpleTrait").Find(&userTraits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userTraits))
			assert.Equal(t, "Cycling", userTraits[0].SimpleTraitName)
			assert.Equal(t, trait, *userTraits[0].SimpleTrait)
		},
		TestName: "Test adding simple trait by name, trait already exists, ignoring case",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitById(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{
				Name:        "Cycling",
				Type:        data.SIMPLE_TRAIT_TYPE_UNDETERMINED,
				IsSensitive: true,
			}
			err := db.Save(&trait).Error
			assert.NoError(t, err)

			err = AddUserSimpleTraitById(db, data.TUserID(1), trait.Id)
			assert.NoError(t, err)

			var traits []data.UserSimpleTrait
			err = db.Where(&data.UserSimpleTrait{UserId: 1}).Preload("SimpleTrait").Find(&traits).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(traits))
			assert.Equal(t, trait.Id, traits[0].SimpleTraitId)
			assert.Equal(t, "Cycling", traits[0].SimpleTraitName)
			assert.Equal(t, data.SIMPLE_TRAIT_TYPE_UNDETERMINED, traits[0].SimpleTraitType)
			assert.True(t, traits[0].SimpleTraitIsSensitive)
			assert.Equal(t, trait, *traits[0].SimpleTrait)
		},
		TestName: "Test adding simple trait by id",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitInvalidId(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			err := AddUserSimpleTraitById(db, data.TUserID(1), data.TSimpleTraitID(1))
			assert.Error(t, err)
		},
		TestName: "Test adding simple trait by id where trait doesn't exist",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitAlreadyHaveMatching(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{Name: "Cycling"}
			err := db.Save(&trait).Error
			assert.NoError(t, err)

			userTrait := data.UserSimpleTrait{UserId: 1, SimpleTraitId: trait.Id}
			err = db.Save(&userTrait).Error
			assert.NoError(t, err)

			err = AddUserSimpleTraitById(db, data.TUserID(1), data.TSimpleTraitID(1))
			assert.Error(t, err)
		},
		TestName: "Test adding simple trait by id where they already have that trait",
	}
	test.RunTestWithDb(thisTest)
}

func TestRemoveUserSimpleTrait(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{Name: "Cycling"}
			err := db.Save(&trait).Error
			assert.NoError(t, err)

			userTrait := data.UserSimpleTrait{UserId: 1, SimpleTraitId: trait.Id}
			err = db.Save(&userTrait).Error
			assert.NoError(t, err)

			var traits []data.UserSimpleTrait
			db.Where(&data.UserSimpleTrait{UserId: 1}).Find(&traits)
			assert.Equal(t, 1, len(traits))

			err = RemoveUserSimpleTrait(db, data.TUserID(1), userTrait.Id)
			assert.NoError(t, err)

			db.Where(&data.UserSimpleTrait{UserId: 1}).Find(&traits)
			assert.Equal(t, 0, len(traits))
		},
		TestName: "Test removing simple trait by id",
	}
	test.RunTestWithDb(thisTest)
}
