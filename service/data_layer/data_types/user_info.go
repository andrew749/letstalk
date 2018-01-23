package data_types

type UserInfo struct {
	Gender    *Gender
	Birthdate *Birthdate
}

/**
 * Validates user info and constructs a safe struct
 */
func CreateUserInfo(
	gender *Gender,
	birthdate *Birthdate,
) (*UserInfo, error) {

	gender, err := CreateGender(gender)

	if err != nil {
		return nil, err
	}

	birthdate, err = CreateBirthdate(birthdate)

	if err != nil {
		return nil, err
	}

	return &UserInfo{gender, birthdate}, nil
}
