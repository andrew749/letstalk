package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type Error interface {
	error
	Code() int
}

type ClientError struct {
	error
}

func (e *ClientError) Code() int { return http.StatusBadRequest }

func NewClientError(msg string, args ...interface{}) *ClientError {
	return &ClientError{errors.New(fmt.Sprintf(msg, args...))}
}


type InternalError struct {
	error
}

func (e *InternalError) Code() int { return http.StatusInternalServerError }

func NewInternalError(msg string, args ...interface{}) *InternalError {
	return &InternalError{errors.New(fmt.Sprintf(msg, args...))}
}
