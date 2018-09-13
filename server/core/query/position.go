package query

import (
	"fmt"
	"time"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
)

func getRole(db *gorm.DB, roleId data.TRoleID) (*data.Role, errs.Error) {
	var role data.Role
	err := db.Where(&data.Role{Id: roleId}).First(&role).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.NewRequestError(fmt.Sprintf("Role with id %d not found", roleId))
		}
		return nil, errs.NewDbError(err)
	}
	return &role, nil
}

// Returns a role with the given name or creates a new one if one doesn't already exist.
func getOrCreateRole(
	db *gorm.DB,
	es *elastic.Client,
	name string,
) (*data.Role, errs.Error) {
	var role data.Role

	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		err := db.Where("UPPER(name) = UPPER(?)", name).First(&role).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				role = data.Role{
					Name:            name,
					IsUserGenerated: true,
				}

				// Add role if it doesn't already exist.
				if err := db.Create(&role).Error; err != nil {
					return err
				}

				go indexRole(es, role)
			} else {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	return &role, nil
}

func getOrganization(db *gorm.DB, orgId data.TOrganizationID) (*data.Organization, errs.Error) {
	var organization data.Organization
	err := db.Where(&data.Organization{Id: orgId}).First(&organization).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.NewRequestError(fmt.Sprintf("Organization with id %d not found", orgId))
		}
		return nil, errs.NewDbError(err)
	}
	return &organization, nil
}

// Returns a organization with the given name or creates a new one if one doesn't already exist.
func getOrCreateOrganization(
	db *gorm.DB,
	es *elastic.Client,
	name string,
) (*data.Organization, errs.Error) {
	var organization data.Organization

	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		err := db.Where("UPPER(name) = UPPER(?)", name).First(&organization).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				organization = data.Organization{
					Name:            name,
					Type:            data.ORGANIZATION_TYPE_UNDETERMINED,
					IsUserGenerated: true,
				}

				// Add organization if it doesn't already exist.
				if err := db.Create(&organization).Error; err != nil {
					return err
				}

				go indexOrganization(es, organization)
			} else {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	return &organization, nil
}

func addUserPosition(
	db *gorm.DB,
	es *elastic.Client,
	userId data.TUserID,
	role data.Role,
	organization data.Organization,
	startDate string,
	endDate *string,
) (*data.UserPosition, errs.Error) {
	userPosition := data.UserPosition{
		UserId:           userId,
		OrganizationId:   organization.Id,
		OrganizationName: organization.Name,
		OrganizationType: organization.Type,
		RoleId:           role.Id,
		RoleName:         role.Name,
		StartDate:        startDate,
		EndDate:          endDate,
	}
	if err := db.Create(&userPosition).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	go indexPositionMultiTrait(es, &userPosition)
	return &userPosition, nil
}

// TODO: Move elsewhere
func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// One of roleId and roleName must be provided.
// One of orgId and orgName must be provided.
func AddUserPosition(
	db *gorm.DB,
	es *elastic.Client,
	userId data.TUserID,
	roleId *data.TRoleID,
	roleName *string,
	organizationId *data.TOrganizationID,
	organizationName *string,
	startDate string,
	endDate *string,
) (*data.UserPosition, errs.Error) {
	if !isValidDate(startDate) {
		return nil, errs.NewRequestError(
			fmt.Sprintf("startDate %s should be in YYYY-MM-DD format", startDate),
		)
	}
	if endDate != nil {
		if !isValidDate(*endDate) {
			return nil, errs.NewRequestError(
				fmt.Sprintf("endDate %s should be in YYYY-MM-DD format", *endDate),
			)
		} else if *endDate <= startDate {
			return nil, errs.NewRequestError(
				fmt.Sprintf("endDate %s should be after startDate %s", *endDate, startDate),
			)
		}
	}

	var (
		role         *data.Role         = nil
		organization *data.Organization = nil
		err          errs.Error         = nil
	)

	if roleId != nil {
		role, err = getRole(db, *roleId)
	} else if roleName != nil {
		role, err = getOrCreateRole(db, es, *roleName)
	} else {
		err = errs.NewRequestError("Must provide either roleId or roleName")
	}
	if err != nil {
		return nil, err
	}

	if organizationId != nil {
		organization, err = getOrganization(db, *organizationId)
	} else if organizationName != nil {
		organization, err = getOrCreateOrganization(db, es, *organizationName)
	} else {
		err = errs.NewRequestError("Must provide either organizationId or organizationName")
	}
	if err != nil {
		return nil, err
	}

	return addUserPosition(db, es, userId, *role, *organization, startDate, endDate)
}

// TODO: Maybe check if exists
func RemoveUserPosition(
	db *gorm.DB,
	userId data.TUserID,
	userPositionId data.TUserPositionID,
) errs.Error {
	toDelete := data.UserPosition{Id: userPositionId, UserId: userId}
	err := db.Delete(&toDelete).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
