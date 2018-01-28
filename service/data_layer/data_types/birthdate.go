package data_types

import (
	"time"
)

type Birthdate time.Time

func CreateBirthdate(birthDate *Birthdate) (*Birthdate, error) {
	if err := validateBirthdate(birthDate); err != nil {
		return nil, err
	}
	return birthDate, nil
}

func validateBirthdate(birthDate *Birthdate) error {
	// FIXME: date checking logic
	return nil
}
