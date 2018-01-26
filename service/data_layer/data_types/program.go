package data_types

import (
	"uwletstalk/service/data_layer/exceptions"
)

type Program string

func CreateProgram(program *Program) (*Program, error) {
	if err := validateProgram(program); err != nil {
		return nil, err
	}
	return program, nil
}

func validateProgram(program *Program) error {
	// TODO: Don't hardcode these
	switch *program {
	case "Software Engineering":
	case "Computer Engineering":
	case "Mechatronics Engineering":
	case "Systems Design Engineering":
		return nil
	}
	return exceptions.CreateClientException("Invalid Program")
}
