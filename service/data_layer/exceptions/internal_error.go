package exceptions

type InternalError struct {
	Exception
}

func (e InternalError) GetErrorCode() int {
	return 500
}

func CreateInternalError() InternalError {
	return InternalError{CreateException("Internal Error")}
}
