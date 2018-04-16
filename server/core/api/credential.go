package api

import (
	"errors"
	"fmt"

	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// Credential Options

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

// Credentials

type CredentialId int

type CredentialPair struct {
	PositionId     data.CredentialPositionId     `json:"positionId"`
	OrganizationId data.CredentialOrganizationId `json:"organizationId"`
}

type CredentialPairWithId struct {
	CredentialPair
	CredentialId CredentialId `json:"credentialId"`
}

type Credential struct {
	PositionId       data.CredentialPositionId     `json:"positionId"`
	PositionName     string                        `json:"positionName"`
	OrganizationId   data.CredentialOrganizationId `json:"organizationId"`
	OrganizationName string                        `json:"organizationName"`
}

type CredentialWithId struct {
	Credential
	CredentialId CredentialId `json:"credentialId"`
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

func resolveCredentialNames(
	credentialPairs []CredentialPairWithId,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) ([]CredentialWithId, errs.Error) {
	credentials := make([]CredentialWithId, len(credentialPairs))
	for i, credentialPair := range credentialPairs {
		org, orgOk := orgMap[credentialPair.OrganizationId]
		if !orgOk {
			return nil, errs.NewClientError(fmt.Sprintf("Invalid organization id %d",
				credentialPair.OrganizationId))
		}

		pos, posOk := posMap[credentialPair.PositionId]
		if !posOk {
			return nil, errs.NewClientError(fmt.Sprintf("Invalid position id %d",
				credentialPair.PositionId))
		}

		credentials[i] = CredentialWithId{
			CredentialId: CredentialId(credentialPair.CredentialId),
			Credential: Credential{
				PositionId:       credentialPair.PositionId,
				PositionName:     pos.Name,
				OrganizationId:   credentialPair.OrganizationId,
				OrganizationName: org.Name,
			},
		}
	}
	return credentials, nil
}

func getUserCredentialsInner(
	db *gorm.DB,
	userId int,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) ([]CredentialWithId, errs.Error) {
	var userCredentials []data.UserCredential
	if err := db.Where("user_id = ?", userId).Find(&userCredentials).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentialPairs := make([]CredentialPairWithId, len(userCredentials))
	for i, userCredential := range userCredentials {
		credentialPairs[i] = CredentialPairWithId{
			CredentialPair: CredentialPair{
				PositionId:     userCredential.PositionId,
				OrganizationId: userCredential.OrganizationId,
			},
			CredentialId: CredentialId(userCredential.ID),
		}
	}
	return resolveCredentialNames(credentialPairs, orgMap, posMap)
}

func GetUserCredentials(db *gorm.DB, userId int) ([]CredentialWithId, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()
	return getUserCredentialsInner(db, userId, orgMap, posMap)
}

func AddUserCredential(
	db *gorm.DB,
	userId int,
	credential CredentialPair,
) (*CredentialId, errs.Error) {
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

	credentialId := CredentialId(newUserCred.ID)
	return &credentialId, nil
}

func RemoveUserCredential(db *gorm.DB, userId int, credentialId CredentialId) errs.Error {
	err := db.Where("id = ? AND user_id = ?", credentialId, userId).Delete(
		data.UserCredential{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

// Credential Requests

// TODO: Get rid of duplicate code by using interfaces

type CredentialRequestId int

type CredentialRequestWithId struct {
	Credential
	CredentialRequestId CredentialRequestId `json:"credentialRequestId"`
}

func getUserCredentialRequestsInner(
	db *gorm.DB,
	userId int,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) ([]CredentialRequestWithId, errs.Error) {
	var userRequests []data.UserCredentialRequest
	if err := db.Where("user_id = ?", userId).Find(&userRequests).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentialPairs := make([]CredentialPairWithId, len(userRequests))
	for i, userRequest := range userRequests {
		credentialPairs[i] = CredentialPairWithId{
			CredentialPair: CredentialPair{
				PositionId:     userRequest.PositionId,
				OrganizationId: userRequest.OrganizationId,
			},
			CredentialId: CredentialId(userRequest.ID),
		}
	}
	credentials, err := resolveCredentialNames(credentialPairs, orgMap, posMap)
	if err != nil {
		return nil, err
	}

	credentialRequests := make([]CredentialRequestWithId, len(credentials))
	for i, credential := range credentials {
		credentialRequests[i] = CredentialRequestWithId{
			Credential:          credential.Credential,
			CredentialRequestId: CredentialRequestId(credential.CredentialId),
		}
	}
	return credentialRequests, nil
}

func GetUserCredentialRequests(
	db *gorm.DB,
	userId int,
) ([]CredentialRequestWithId, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()
	return getUserCredentialRequestsInner(db, userId, orgMap, posMap)
}

func AddUserCredentialRequest(
	db *gorm.DB,
	userId int,
	credential CredentialPair,
) (*CredentialRequestId, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()

	if err := validateCredential(credential, orgMap, posMap); err != nil {
		return nil, errs.NewClientError(err.Error())
	}

	userRequests, err := getUserCredentialRequestsInner(db, userId, orgMap, posMap)
	if err != nil {
		return nil, err
	}

	for _, userRequest := range userRequests {
		if userRequest.OrganizationId == credential.OrganizationId &&
			userRequest.PositionId == credential.PositionId {
			return nil, errs.NewClientError("You have already requested this credential")
		}
	}

	newReq := data.UserCredentialRequest{
		UserId:         userId,
		PositionId:     credential.PositionId,
		OrganizationId: credential.OrganizationId,
	}
	if err := db.Create(&newReq).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentialRequestId := CredentialRequestId(newReq.ID)
	return &credentialRequestId, nil
}

func RemoveUserCredentialRequest(
	db *gorm.DB,
	userId int,
	credentialRequestId CredentialRequestId,
) errs.Error {
	err := db.Where("id = ? AND user_id = ?", credentialRequestId, userId).Delete(
		data.UserCredentialRequest{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
