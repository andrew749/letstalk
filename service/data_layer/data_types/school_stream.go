package data_types

type SchoolStream string

type SchoolStreamError struct {
	ErrorCode int
}

func CreateSchoolStream(schoolStream *SchoolStream) (*SchoolStream, error) {
	if err := validateSchoolStream(schoolStream); err != nil {
		return nil, err
	}
	return schoolStream, nil
}

func validateSchoolStream(schoolStream *SchoolStream) error {
	switch *schoolStream {
	case "4 Stream":
	case "8 Stream":
		return nil
	}
	return SchoolStreamError{0}
}

func (sse SchoolStreamError) Error() string {
	switch sse.ErrorCode {
	case 0:
		return "Bad School Term"
	}
	return "Unknown Error"
}
