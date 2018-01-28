package data_types

import (
	"uwletstalk/service/data_layer/exceptions"
)

type SchoolStream string

func CreateSchoolStream(schoolStream *SchoolStream) (*SchoolStream, error) {
	if err := validateSchoolStream(schoolStream); err != nil {
		return nil, err
	}
	return schoolStream, nil
}

func validateSchoolStream(schoolStream *SchoolStream) error {
	// TODO: remove this hardcoding possibly
	switch *schoolStream {
	case "4 Stream":
	case "8 Stream":
		return nil
	}
	return exceptions.CreateClientException("Unknown Stream")
}
