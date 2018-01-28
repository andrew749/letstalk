package data_types

import (
	"uwletstalk/service/data_layer/exceptions"
)

type Gender int

const (
	MALE   Gender = 0
	FEMALE Gender = 1
	OTHER  Gender = 10
)

func CreateGender(genderCode *Gender) (*Gender, error) {
	if err := validateGender(genderCode); err != nil {
		return nil, err
	}
	return genderCode, nil
}

func validateGender(gender *Gender) error {
	switch *gender {
	case MALE:
		fallthrough
	case FEMALE:
		fallthrough
	case OTHER:
		return nil
	}
	return exceptions.CreateClientException("Unknown Gender")
}
