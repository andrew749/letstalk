package exceptions

/**
 * ABC for exceptions. No class should instantiate an exceptions
 * on it's own; rather it should create a more specific subclass.
 */
type Exception struct {
	Message string
}

type IException interface {
	GetErrorCode() int
}

func (e Exception) Error() string {
	return e.Message
}

func CreateException(
	message string,
) Exception {
	return Exception{message}
}
