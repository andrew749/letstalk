package query

import (
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

type CredentialPair struct {
	PositionId     data.CredentialPositionId     `json:"positionId"`
	OrganizationId data.CredentialOrganizationId `json:"organizationId"`
}

type CredentialPairWithId struct {
	CredentialPair
	CredentialId CredentialId `json:"credentialId"`
}

type CredentialId int

type CredentialStrategy interface {
	GetCredentialsForUser() ([]CredentialPairWithId, errs.Error)
	CreateCredentialForUser(credentialPair CredentialPair) (*CredentialId, errs.Error)
	DeleteCredentialForUser(credentialId CredentialId) errs.Error
	NewDuplicateError() errs.Error
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
) errs.Error {
	org, orgOk := orgMap[credential.OrganizationId]
	if !orgOk {
		return errs.NewClientError(fmt.Sprintf("Invalid organization id %d", credential.OrganizationId))
	}

	pos, posOk := posMap[credential.PositionId]
	if !posOk {
		return errs.NewClientError(fmt.Sprintf("Invalid position id %d", credential.PositionId))
	}

	for _, pair := range validPairs {
		if pair.PositionType == pos.Type && pair.OrganizationType == org.Type {
			return nil
		}
	}
	return errs.NewClientError(fmt.Sprintf("Invalid organization type, position type pair (%d, %d)",
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

type UserCredentialStrategy struct {
	Db     *gorm.DB
	UserId int
}

func (s UserCredentialStrategy) GetCredentialsForUser() ([]CredentialPairWithId, errs.Error) {
	var userCredentials []data.UserCredential
	if err := s.Db.Where("user_id = ?", s.UserId).Find(&userCredentials).Error; err != nil {
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

	return credentialPairs, nil
}

func (s UserCredentialStrategy) CreateCredentialForUser(
	credentialPair CredentialPair,
) (*CredentialId, errs.Error) {
	newUserCred := data.UserCredential{
		UserId:         s.UserId,
		PositionId:     credentialPair.PositionId,
		OrganizationId: credentialPair.OrganizationId,
	}
	if err := s.Db.Create(&newUserCred).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentialId := CredentialId(newUserCred.ID)
	return &credentialId, nil
}

func (s UserCredentialStrategy) DeleteCredentialForUser(
	credentialId CredentialId,
) errs.Error {
	err := s.Db.Where("id = ? AND user_id = ?", credentialId, s.UserId).Delete(
		data.UserCredential{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func (s UserCredentialStrategy) NewDuplicateError() errs.Error {
	return errs.NewClientError("You already have this credential")
}

func getCredentialsWithStrategyInner(
	s CredentialStrategy,
	orgMap map[data.CredentialOrganizationId]data.CredentialOrganization,
	posMap map[data.CredentialPositionId]data.CredentialPosition,
) ([]CredentialWithId, errs.Error) {
	credentialPairs, err := s.GetCredentialsForUser()
	if err != nil {
		return nil, err
	}
	return resolveCredentialNames(credentialPairs, orgMap, posMap)
}

func GetCredentialsWithStrategy(s CredentialStrategy) ([]CredentialWithId, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()
	return getCredentialsWithStrategyInner(s, orgMap, posMap)
}

func AddCredentialWithStrategy(
	s CredentialStrategy,
	credential CredentialPair,
) (*CredentialId, errs.Error) {
	orgMap := data.BuildOrganizationIdIndex()
	posMap := data.BuildPositionIdIndex()

	if err := validateCredential(credential, orgMap, posMap); err != nil {
		return nil, errs.NewClientError(err.Error())
	}

	existingCreds, err := getCredentialsWithStrategyInner(s, orgMap, posMap)
	if err != nil {
		return nil, err
	}

	for _, existingCred := range existingCreds {
		if existingCred.OrganizationId == credential.OrganizationId &&
			existingCred.PositionId == credential.PositionId {
			return nil, s.NewDuplicateError()
		}
	}

	return s.CreateCredentialForUser(credential)
}

// Credential Requests

type CredentialRequestId int

type CredentialRequestWithId struct {
	Credential
	CredentialRequestId CredentialRequestId `json:"credentialRequestId"`
}

type UserCredentialRequestStrategy struct {
	Db     *gorm.DB
	UserId int
}

func (s UserCredentialRequestStrategy) GetCredentialsForUser() (
	[]CredentialPairWithId,
	errs.Error,
) {
	var userRequests []data.UserCredentialRequest
	if err := s.Db.Where("user_id = ?", s.UserId).Find(&userRequests).Error; err != nil {
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
	return credentialPairs, nil
}

func (s UserCredentialRequestStrategy) CreateCredentialForUser(
	credentialPair CredentialPair,
) (*CredentialId, errs.Error) {
	newReq := data.UserCredentialRequest{
		UserId:         s.UserId,
		PositionId:     credentialPair.PositionId,
		OrganizationId: credentialPair.OrganizationId,
	}
	if err := s.Db.Create(&newReq).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentialId := CredentialId(newReq.ID)
	return &credentialId, nil
}

func (s UserCredentialRequestStrategy) DeleteCredentialForUser(
	credentialId CredentialId,
) errs.Error {
	err := s.Db.Where("id = ? AND user_id = ?", credentialId, s.UserId).Delete(
		data.UserCredentialRequest{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func (s UserCredentialRequestStrategy) NewDuplicateError() errs.Error {
	return errs.NewClientError("You have already requested this credential")
}
