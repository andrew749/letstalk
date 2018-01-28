package data_types

type SchoolInfo struct {
	Program *Program
	Stream  *SchoolStream
}

func CreateSchoolInfo(program *Program, stream *SchoolStream) (*SchoolInfo, error) {
	program, err := CreateProgram(program)
	if err != nil {
		return nil, err
	}

	stream, err = CreateSchoolStream(stream)
	if err != nil {
		return nil, err
	}

	return &SchoolInfo{program, stream}, nil
}
