package errs

import (
	"net/http"
)

type BadRequest struct{ IError }        // 400
type UnauthorizedError struct{ IError } // 401
type ForbiddenError struct{ IError }    // 403
type NotFoundError struct{ IError }     // 404

func (e *BadRequest) GetHTTPCode() int        { return http.StatusBadRequest }
func (e *UnauthorizedError) GetHTTPCode() int { return http.StatusUnauthorized }
func (e *ForbiddenError) GetHTTPCode() int    { return http.StatusForbidden }
func (e *NotFoundError) GetHTTPCode() int     { return http.StatusNotFound }

func NewRequestError(msg string, args ...interface{}) IError {
	return &BadRequest{NewBaseError(msg, args...)}
}
func NewUnauthorizedError(msg string, args ...interface{}) IError {
	return &UnauthorizedError{NewBaseError(msg, args...)}
}
func NewForbiddenError(msg string, args ...interface{}) IError {
	return &ForbiddenError{NewBaseError(msg, args...)}
}
func NewNotFoundError(msg string, args ...interface{}) IError {
	return &NotFoundError{NewBaseError(msg, args...)}
}

type InternalError struct{ IError }
type DatabaseError struct{ IError }
type ElasticsearchError struct{ IError }

func (e *InternalError) GetHTTPCode() int { return http.StatusInternalServerError }

func NewInternalError(msg string, args ...interface{}) IError {
	return &InternalError{
		NewBaseError(msg, args...),
	}
}

func NewDbError(err error) IError {
	return &DatabaseError{NewInternalError("Encountered database error: %s", err)}
}

func NewEsError(err error) IError {
	return &ElasticsearchError{NewInternalError("Encountered elasticsearch error: %s", err)}
}

type InvalidPasswordError struct {
	IError
}

func InvalidPassError() Error {
	return &InvalidPasswordError{
		NewRequestError("Invalid Password. Try again."),
	}
}
