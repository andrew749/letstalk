package data_types

import (
  "./program"
  "./school_stream"
)

type SchoolInfo struct {
	Program Program
	Stream  SchoolStream
}

func CreateSchoolInfo(program Program, stream SchoolStream) SchoolInfo, error {
  program, err := CreateProgram(program)
  if err != nil {
    return err
  }

  stream, err := CreateSchoolStream(stream)
  if err != nil {
    return err
  }

  return SchoolInfo{program, stfream}
}
