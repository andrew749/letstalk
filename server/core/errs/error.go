package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type Error interface {
	error
	GetHTTPCode() int
}

type ClientError struct {
	error
}

func (e *ClientError) GetHTTPCode() int { return http.StatusBadRequest }

func NewClientError(msg string, args ...interface{}) Error {
	return &ClientError{errors.New(fmt.Sprintf(msg, args...))}
}

type InternalError struct {
	error
}

func (e *InternalError) GetHTTPCode() int { return http.StatusInternalServerError }

func NewInternalError(msg string, args ...interface{}) Error {
	return &InternalError{errors.New(fmt.Sprintf(msg, args...))}
}

func NewDbError(err error) Error {
	return NewInternalError("encountered database error: %s", err)
}
