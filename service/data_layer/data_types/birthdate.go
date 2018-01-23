package data_types

import (
	"time"
)

type Birthdate time.Time

type BirthdateError struct {
	ErrorCode int
}

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

func (be BirthdateError) Error() string {
	switch be.ErrorCode {
	case 0:
		return "Malformed date"
	case 1:
		return "Not old enough"
	}
	return "Unknown Error"
}
