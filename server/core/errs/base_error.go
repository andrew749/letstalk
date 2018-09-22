package errs

import (
	"fmt"

	"github.com/pkg/errors"
)

type IError interface {
	error
	GetHTTPCode() int
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

func (e *BaseError) GetHTTPCode() int { panic("Abstract Error") }

func (e *BaseError) GetExtraData() map[string]interface{} {
	return e.ExtraData
}

func NewBaseError(msg string, args ...interface{}) *BaseError {
	extraData := make(map[string]interface{})
	// add stack trace context information
	return &BaseError{errors.New(fmt.Sprintf(msg, args...)), extraData}
}
