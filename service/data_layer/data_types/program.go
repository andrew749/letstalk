package data_types

type Program string

type ProgramError struct {
	ErrorCode int
}

func CreateProgram(program *Program) (*Program, error) {
	if err := validateProgram(program); err != nil {
		return nil, err
	}
	return program, nil
}

func validateProgram(program *Program) error {
	switch *program {
	case "Software Engineering":
	case "Computer Engineering":
	case "Mechatronics Engineering":
	case "Systems Design Engineering":
		return nil
	}
	return ProgramError{0}
}

func (pe ProgramError) Error() string {
	switch pe.ErrorCode {
	case 0:
		return "Bad Program"
	}
	return "Unknown Error"
}
