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

type CredentialPair struct {
	PositionId     data.CredentialPositionId     `json:"positionId"`
	OrganizationId data.CredentialOrganizationId `json:"organizationId"`
}

type Credential struct {
	PositionId       data.CredentialPositionId     `json:"positionId"`
	PositionName     string                        `json:"positionName"`
	OrganizationId   data.CredentialOrganizationId `json:"organizationId"`
	OrganizationName string                        `json:"organizationName"`
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

func validateCredential(
	credential CredentialPair,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) error {
	org, orgOk := orgMap[credential.OrganizationId]
	if !orgOk {
		return errors.New(fmt.Sprintf("Invalid organization id %d",
			credential.OrganizationId))
	}

	pos, posOk := posMap[credential.PositionId]
	if !posOk {
		return errors.New(fmt.Sprintf("Invalid position id %d", credential.PositionId))
	}

	for _, pair := range validPairs {
		if pair.PositionType == pos.Type && pair.OrganizationType == org.Type {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid organization type, position type pair (%d, %d)",
		org.Type, pos.Type))
}

func getUserCredentialsInner(
	db *gorm.DB,
	userId int,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) ([]Credential, errs.Error) {
	var userCredentials []data.UserCredential
	if err := db.Where("user_id = ?", userId).Find(&userCredentials).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentials := make([]Credential, len(userCredentials))
	for i, userCredential := range userCredentials {
		org, orgOk := orgMap[userCredential.OrganizationId]
		if !orgOk {
			return nil, errs.NewClientError(fmt.Sprintf("Invalid organization id %d",
				userCredential.OrganizationId))
		}

		pos, posOk := posMap[userCredential.PositionId]
		if !posOk {
			return nil, errs.NewClientError(fmt.Sprintf("Invalid position id %d",
				userCredential.PositionId))
		}

		credentials[i] = Credential{
			PositionId:       userCredential.PositionId,
			PositionName:     pos.Name,
			OrganizationId:   userCredential.OrganizationId,
			OrganizationName: org.Name,
		}
	}
	return credentials, nil
}

func GetUserCredentials(db *gorm.DB, userId int) ([]Credential, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()
	return getUserCredentialsInner(db, userId, orgMap, posMap)
}

func AddUserCredential(db *gorm.DB, userId int, credential CredentialPair) (*uint, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()

	if err := validateCredential(credential, orgMap, posMap); err != nil {
		return nil, errs.NewClientError(err.Error())
	}

	userCredentials, err := getUserCredentialsInner(db, userId, orgMap, posMap)
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
