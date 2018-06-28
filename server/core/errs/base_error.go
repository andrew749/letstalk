package errs

import (
	"errors"
	"fmt"
)

type IError interface {
	error
	GetHTTPCode() int
	GetExtraData() map[string]interface{} // key value attributes associated with error
}

type Error = IError

type BaseError struct {
	error
	ExtraData map[string]interface{}
}

func (e *BaseError) SetExtra(key string, value interface{}) {
	e.ExtraData[key] = value
}

func (e *BaseError) Error() string {
	return ""
}

func (e *BaseError) GetHTTPCode() int { panic("Abstract Error") }

func (e *BaseError) GetExtraData() map[string]interface{} {
	return e.ExtraData
}

func NewBaseError(msg string, args ...interface{}) *BaseError {
	extraData := make(map[string]interface{})
	return &BaseError{errors.New(fmt.Sprintf(msg, args...)), extraData}
}
