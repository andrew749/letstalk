package data_types

type Gender int

type GenderError struct {
	ErrorCode int
}

const (
	MALE   Gender = 0
	FEMALE Gender = 1
  OTHER Gender = 10
)

func CreateGender(genderCode int) Gender, error {
    if err := validateGender(genderCode) != nil {
      return nil, err
    }
    return genderCode, nil
}

func validateGender(gender Gender) error {
	switch gender {
	case MALE:
	case FEMALE:
  case OTHER:
		return nil
	}
	return GenderError{0}
}

func (ge GenderError) Error() string {
	switch ge.ErrorCode {
	case 0:
		return "Invalid Gender"
	}
	return "Unknown Error"
}
