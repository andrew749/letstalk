package data_types

import (
  "gender"
)

type UserInfo struct {
	Gender    Gender
	Birthdate Time
}

/**
 * Validates user info and constructs a safe struct
 */
func CreateUserInfo(
	gender Gender
	birthdate Time
) UserInfo, error {

  gender, err := CreateGender(genderCode)

  if err != nil {
    return nil, err
  }

  birthdate, err := CreateBirthdate(birthdate)

  if err != nil {
    return nil, err
  }

  return UserInfo{gender, birthdate}

}
