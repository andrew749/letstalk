package api

import (
	"errors"
	"fmt"

	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

type ValidCredentialPair struct {
	PositionType     data.CredentialPositionType     `json:"positionType"`
	OrganizationType data.CredentialOrganizationType `json:"organizationType"`
}

type CredentialOptions struct {
	ValidPairs    []ValidCredentialPair         `json:"validPairs"`
	Organizations []data.CredentialOrganization `json:"organizations"`
	Positions     []data.CredentialPosition     `json:"positions"`
}

var validPairs = []ValidCredentialPair{
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_COOP,
		data.CREDENTIAL_ORGANIZATION_TYPE_COOP,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_CLUB,
		data.CREDENTIAL_ORGANIZATION_TYPE_CLUB,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_SPORT,
		data.CREDENTIAL_ORGANIZATION_TYPE_SPORT,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_COHORT,
		data.CREDENTIAL_ORGANIZATION_TYPE_COHORT,
	},
}

type Credential struct {
	PositionId     data.CredentialPositionId     `json:"positionId"`
	OrganizationId data.CredentialOrganizationId `json:"organizationId"`
}

// Returns a struct contain all info required to generate all possible credential options, where
// a credential consists of a position and an organization.
func GetCredentialOptions() CredentialOptions {
	// TODO: Could cache the results of BuildOrganizations and BuildPositions
	return CredentialOptions{
		ValidPairs:    validPairs,
		Organizations: data.BuildOrganizations(),
		Positions:     data.BuildPositions(),
	}
}

func validateCredential(credential Credential) error {
	orgInverseTypeMap := data.BuildInverseOrganizationTypeMap()
	posInverseTypeMap := data.BuildInversePositionTypeMap()

	orgType, orgOk := orgInverseTypeMap[credential.OrganizationId]
	if !orgOk {
		return errors.New(fmt.Sprintf("Invalid organization id %d",
			credential.OrganizationId))
	}

	posType, posOk := posInverseTypeMap[credential.PositionId]
	if !posOk {
		return errors.New(fmt.Sprintf("Invalid position id %d", credential.PositionId))
	}

	for _, pair := range validPairs {
		if pair.PositionType == posType && pair.OrganizationType == orgType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid organization type, position type pair (%d, %d)",
		orgType, posType))
}

func GetUserCredentials(db *gorm.DB, userId int) ([]data.UserCredential, errs.Error) {
	var userCredentials []data.UserCredential
	if err := db.Where("user_id = ?", userId).Find(&userCredentials).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	return userCredentials, nil
}

func AddUserCredential(db *gorm.DB, userId int, credential Credential) (*uint, errs.Error) {
	if err := validateCredential(credential); err != nil {
		return nil, errs.NewClientError(err.Error())
	}

	userCredentials, err := GetUserCredentials(db, userId)
	if err != nil {
		return nil, err
	}

	for _, userCred := range userCredentials {
		if userCred.OrganizationId == credential.OrganizationId &&
			userCred.PositionId == credential.PositionId {
			return nil, errs.NewClientError("You already have this credential")
		}
	}

	newUserCred := data.UserCredential{
		UserId:         userId,
		PositionId:     credential.PositionId,
		OrganizationId: credential.OrganizationId,
	}
	if err := db.Create(&newUserCred).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	return &newUserCred.ID, nil
}

func RemoveUserCredential(db *gorm.DB, userId int, credentialId uint) errs.Error {
	err := db.Where("id = ? AND user_id = ?", credentialId, userId).Delete(
		data.UserCredential{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
