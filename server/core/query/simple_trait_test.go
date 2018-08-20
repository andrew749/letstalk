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
			err := AddUserSimpleTraitByName(db, data.TUserID(1), "Cycling ")
			assert.Nil(t, err)

			var traits []data.UserSimpleTrait
			db.Where(&data.UserSimpleTrait{UserId: 1}).Preload("SimpleTrait").Find(&traits)
			assert.Equal(t, 1, len(traits))
			assert.Equal(t, "Cycling", traits[0].SimpleTraitName)
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
			db.Save(&trait)
			var traits []data.SimpleTrait
			db.Find(&traits)
			assert.Equal(t, 1, len(traits))

			err := AddUserSimpleTraitByName(db, data.TUserID(1), "Cycling ")
			assert.Nil(t, err)

			db.Find(&traits)
			assert.Equal(t, 1, len(traits))

			var userTraits []data.UserSimpleTrait
			db.Where(&data.UserSimpleTrait{UserId: 1}).Preload("SimpleTrait").Find(&userTraits)
			assert.Equal(t, 1, len(userTraits))
			assert.Equal(t, "Cycling", userTraits[0].SimpleTraitName)
			assert.NotNil(t, userTraits[0].SimpleTrait)
			assert.Equal(t, "Cycling", userTraits[0].SimpleTrait.Name)
			assert.False(t, userTraits[0].SimpleTrait.IsUserGenerated)
		},
		TestName: "Test adding simple trait by name, trait alreadt exists",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitById(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{Name: "Cycling"}
			db.Save(&trait)

			err := AddUserSimpleTraitById(db, data.TUserID(1), trait.Id)
			assert.Nil(t, err)

			var traits []data.UserSimpleTrait
			db.Where(&data.UserSimpleTrait{UserId: 1}).Preload("SimpleTrait").Find(&traits)
			assert.Equal(t, 1, len(traits))
			assert.Equal(t, trait.Id, traits[0].SimpleTraitId)
			assert.Equal(t, "Cycling", traits[0].SimpleTraitName)
			assert.NotNil(t, traits[0].SimpleTrait)
			assert.Equal(t, "Cycling", traits[0].SimpleTrait.Name)
		},
		TestName: "Test adding simple trait by id",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitInvalidId(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			err := AddUserSimpleTraitById(db, data.TUserID(1), data.TSimpleTraitID(1))
			assert.NotNil(t, err)
		},
		TestName: "Test adding simple trait by id where trait doesn't exist",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserSimpleTraitAlreadyHaveMatching(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{Name: "Cycling"}
			db.Save(&trait)
			userTrait := data.UserSimpleTrait{UserId: 1, SimpleTraitId: trait.Id}
			db.Save(&userTrait)

			err := AddUserSimpleTraitById(db, data.TUserID(1), data.TSimpleTraitID(1))
			assert.NotNil(t, err)
		},
		TestName: "Test adding simple trait by id where they already have that trait",
	}
	test.RunTestWithDb(thisTest)
}

func TestRemoveUserSimpleTrait(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			trait := data.SimpleTrait{Name: "Cycling"}
			db.Save(&trait)
			userTrait := data.UserSimpleTrait{UserId: 1, SimpleTraitId: trait.Id}
			db.Save(&userTrait)

			var traits []data.UserSimpleTrait
			db.Where(&data.UserSimpleTrait{UserId: 1}).Find(&traits)
			assert.Equal(t, 1, len(traits))

			err := RemoveUserSimpleTrait(db, data.TUserID(1), trait.Id)
			assert.Nil(t, err)

			db.Where(&data.UserSimpleTrait{UserId: 1}).Find(&traits)
			assert.Equal(t, 0, len(traits))
		},
		TestName: "Test removing simple trait by id",
	}
	test.RunTestWithDb(thisTest)
}
