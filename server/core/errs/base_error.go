package errs

import (
	"fmt"

	"github.com/pkg/errors"
)

type IError interface {
	error
	GetHTTPCode() int
	StackTrace() errors.StackTrace
	GetExtraData() map[string]interface{} // key value attributes associated with error
	VerboseError() string
}

type Error IError

type BaseError struct {
	err       error
	ExtraData map[string]interface{}
}

func (e *BaseError) SetExtra(key string, value interface{}) {
	e.ExtraData[key] = value
}

func (e *BaseError) Error() string {
	return e.err.Error()
}

// VerboseError Provide stack trace information to ease debugging.
func (e *BaseError) VerboseError() string {
	return fmt.Sprintf("%+v", e.err)
}

func (e *BaseError) StackTrace() errors.StackTrace {
	return e.err.(interface {
		StackTrace() errors.StackTrace
	}).StackTrace()
}

func (e *BaseError) GetHTTPCode() int { panic("Abstract Error") }

func (e *BaseError) GetExtraData() map[string]interface{} {
	return e.ExtraData
}

func NewBaseError(msg string, args ...interface{}) *BaseError {
	extraData := make(map[string]interface{})
	return &BaseError{
		errors.Errorf(msg, args...),
		extraData,
	}
}
