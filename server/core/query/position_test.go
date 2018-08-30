package query

import (
	"letstalk/server/core/test"
	"letstalk/server/data"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddUserPositionByIds(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			role := data.Role{
				Name:            "Software Engineering Intern",
				IsUserGenerated: false,
			}
			err := db.Save(&role).Error
			assert.NoError(t, err)

			org := data.Organization{
				Name:            "Facebook",
				Type:            data.ORGANIZATION_TYPE_COMPANY,
				IsUserGenerated: false,
			}
			err = db.Save(&org).Error
			assert.NoError(t, err)

			startDate := "2018-01-01"
			err = AddUserPosition(db, nil, data.TUserID(1), &role.Id, nil, &org.Id, nil, startDate, nil)
			assert.NoError(t, err)

			var userPositions []data.UserPosition
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("Role").Preload("Organization").Find(&userPositions).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userPositions))
			assert.Equal(t, org.Id, userPositions[0].OrganizationId)
			assert.Equal(t, org.Name, userPositions[0].OrganizationName)
			assert.Equal(t, org.Type, userPositions[0].OrganizationType)
			assert.Equal(t, role.Id, userPositions[0].RoleId)
			assert.Equal(t, role.Name, userPositions[0].RoleName)
			assert.Equal(t, startDate, userPositions[0].StartDate)
			assert.Nil(t, userPositions[0].EndDate)
			assert.Equal(t, role, *userPositions[0].Role)
			assert.Equal(t, org, *userPositions[0].Organization)
		},
		TestName: "Test adding user position by ids",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionByNamesAlreadyExist(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			role := data.Role{
				Name:            "Software Engineering Intern",
				IsUserGenerated: false,
			}
			err := db.Save(&role).Error
			assert.NoError(t, err)

			org := data.Organization{
				Name:            "Facebook",
				Type:            data.ORGANIZATION_TYPE_COMPANY,
				IsUserGenerated: false,
			}
			err = db.Save(&org).Error
			assert.NoError(t, err)

			startDate := "2018-01-01"
			endDate := "2018-04-01"
			err = AddUserPosition(
				db, nil, data.TUserID(1), nil, &role.Name, nil, &org.Name, startDate, &endDate,
			)
			assert.NoError(t, err)

			var userPositions []data.UserPosition
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("Role").Preload("Organization").Find(&userPositions).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userPositions))
			assert.Equal(t, org.Id, userPositions[0].OrganizationId)
			assert.Equal(t, org.Name, userPositions[0].OrganizationName)
			assert.Equal(t, org.Type, userPositions[0].OrganizationType)
			assert.Equal(t, role.Id, userPositions[0].RoleId)
			assert.Equal(t, role.Name, userPositions[0].RoleName)
			assert.Equal(t, startDate, userPositions[0].StartDate)
			assert.Equal(t, endDate, *userPositions[0].EndDate)
			assert.Equal(t, role, *userPositions[0].Role)
			assert.Equal(t, org, *userPositions[0].Organization)
		},
		TestName: "Test adding user position by names which already exist",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionByNamesAlreadyExistIgnoreCase(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			role := data.Role{
				Name:            "Software Engineering Intern",
				IsUserGenerated: false,
			}
			err := db.Save(&role).Error
			assert.NoError(t, err)

			org := data.Organization{
				Name:            "Facebook",
				Type:            data.ORGANIZATION_TYPE_COMPANY,
				IsUserGenerated: false,
			}
			err = db.Save(&org).Error
			assert.NoError(t, err)

			roleName := strings.ToLower(role.Name)
			orgName := strings.ToLower(org.Name)
			startDate := "2018-01-01"
			endDate := "2018-04-01"
			err = AddUserPosition(
				db, nil, data.TUserID(1), nil, &roleName, nil, &orgName, startDate, &endDate,
			)
			assert.NoError(t, err)

			var userPositions []data.UserPosition
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("Role").Preload("Organization").Find(&userPositions).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userPositions))
			assert.Equal(t, org.Id, userPositions[0].OrganizationId)
			assert.Equal(t, org.Name, userPositions[0].OrganizationName)
			assert.Equal(t, org.Type, userPositions[0].OrganizationType)
			assert.Equal(t, role.Id, userPositions[0].RoleId)
			assert.Equal(t, role.Name, userPositions[0].RoleName)
			assert.Equal(t, startDate, userPositions[0].StartDate)
			assert.Equal(t, endDate, *userPositions[0].EndDate)
			assert.Equal(t, role, *userPositions[0].Role)
			assert.Equal(t, org, *userPositions[0].Organization)
		},
		TestName: "Test adding user position by names which already exist while ignoring case",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionByNameNotExists(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			roleName := "Software Engineering Intern"
			orgName := "Facebook"
			startDate := "2018-01-01"
			endDate := "2018-04-01"
			var err error
			err = AddUserPosition(
				db, nil, data.TUserID(1), nil, &roleName, nil, &orgName, startDate, &endDate,
			)
			assert.NoError(t, err)

			var userPositions []data.UserPosition
			err = db.Where(
				&data.UserSimpleTrait{UserId: 1},
			).Preload("Role").Preload("Organization").Find(&userPositions).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userPositions))
			assert.Equal(t, orgName, userPositions[0].OrganizationName)
			assert.Equal(t, data.ORGANIZATION_TYPE_UNDETERMINED, userPositions[0].OrganizationType)
			assert.Equal(t, roleName, userPositions[0].RoleName)
			assert.Equal(t, startDate, userPositions[0].StartDate)
			assert.Equal(t, endDate, *userPositions[0].EndDate)
			assert.Equal(t, roleName, userPositions[0].Role.Name)
			assert.True(t, userPositions[0].Role.IsUserGenerated)
			assert.Equal(t, orgName, userPositions[0].Organization.Name)
			assert.Equal(t, data.ORGANIZATION_TYPE_UNDETERMINED, userPositions[0].Organization.Type)
			assert.True(t, userPositions[0].Organization.IsUserGenerated)
		},
		TestName: "Test adding user position by names which don't already exist",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionMissingRole(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			orgName := "Facebook"
			err := AddUserPosition(db, nil, data.TUserID(1), nil, nil, nil, &orgName, "2018-01-01", nil)
			assert.Error(t, err)
			assert.Equal(t, "Must provide either roleId or roleName", err.Error())
		},
		TestName: "Test add user position no role provided",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionMissingOrg(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			roleName := "Software Engineering Intern"
			err := AddUserPosition(db, nil, data.TUserID(1), nil, &roleName, nil, nil, "2018-01-01", nil)
			assert.Error(t, err)
			assert.Equal(t, "Must provide either organizationId or organizationName", err.Error())
		},
		TestName: "Test add user position no org provided",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserPositionInvalidStartDate(t *testing.T) {
	err := AddUserPosition(
		nil, nil, data.TUserID(1), nil, nil, nil, nil, "Monday, June 2nd, 2018", nil)
	assert.Error(t, err)
	assert.Equal(t, "startDate Monday, June 2nd, 2018 should be in YYYY-MM-DD format", err.Error())
}

func TestAddUserPositionInvalidEndDate(t *testing.T) {
	endDate := "Monday, June 2nd, 2018"
	err := AddUserPosition(nil, nil, data.TUserID(1), nil, nil, nil, nil, "2018-01-01", &endDate)
	assert.Error(t, err)
	assert.Equal(t, "endDate Monday, June 2nd, 2018 should be in YYYY-MM-DD format", err.Error())
}

func TestRemoveUserPosition(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			role := data.Role{
				Name:            "Software Engineering Intern",
				IsUserGenerated: false,
			}
			err := db.Save(&role).Error
			assert.NoError(t, err)

			org := data.Organization{
				Name:            "Facebook",
				Type:            data.ORGANIZATION_TYPE_COMPANY,
				IsUserGenerated: false,
			}
			err = db.Save(&org).Error
			assert.NoError(t, err)

			userPosition := data.UserPosition{UserId: 1, RoleId: role.Id, OrganizationId: org.Id}
			err = db.Save(&userPosition).Error
			assert.NoError(t, err)

			var positions []data.UserPosition
			db.Where(&data.UserPosition{UserId: 1}).Find(&positions)
			assert.Equal(t, 1, len(positions))

			err = RemoveUserPosition(db, data.TUserID(1), userPosition.Id)
			assert.NoError(t, err)

			db.Where(&data.UserPosition{UserId: 1}).Find(&positions)
			assert.Equal(t, 0, len(positions))
		},
		TestName: "Test removing user position by id",
	}
	test.RunTestWithDb(thisTest)
}
