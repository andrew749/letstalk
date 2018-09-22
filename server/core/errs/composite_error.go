package errs

import (
	"bytes"
	"fmt"
)

/**
 * Error class to hold multiple errors.
 * Useful when operations are applied on an array and multiple operations
 * can fail.
 */
type CompositeError struct {
	errors   []error
	HttpCode *int
}

func CreateCompositeError() *CompositeError {
	ok := 200
	return &CompositeError{
		make([]error, 0),
		&ok,
	}
}

/**
 * If the composite error is nil and the error to add is not, create a new one
 * If the error is nil and the composite error is nil, return nil
 * Otherwide do the logical append
 */
func AppendNullableError(ce *CompositeError, err error) *CompositeError {
	var e *CompositeError = ce
	if err != nil {
		if ce == nil {
			e = CreateCompositeError()
		}
		e.AddError(err)
	}

	return e
}

func (ce *CompositeError) GetHTTPCode() int {
	return *ce.HttpCode
}

func (ce *CompositeError) AddError(err error) {
	ce.errors = append(ce.errors, err)
}

const ERROR_FORMAT_STRING = "Error %d:\n%s"

func (ce *CompositeError) VerboseError() string {
	var buffer bytes.Buffer
	for i, err := range ce.errors {
		var errMsg string
		switch err.(type) {
		case *BaseError:
			errSpecific := err.(*BaseError)
			errMsg = errSpecific.VerboseError()
		default:
			errMsg = err.Error()
		}
		buffer.WriteString(fmt.Sprintf(ERROR_FORMAT_STRING, i, errMsg))
	}
	return buffer.String()
}

func (ce CompositeError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("Errors: [")
	for _, err := range ce.errors {
		buffer.WriteString("(")
		buffer.WriteString(err.Error())
		buffer.WriteString(")")
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (ce *CompositeError) GetExtraData() map[string]interface{} {
	return map[string]interface{}{}
}
