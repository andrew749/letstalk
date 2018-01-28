package exceptions

type ClientException struct {
	Exception
}

func (e ClientException) GetErrorCode() int {
	return 400
}

func CreateClientException(
	message string,
) ClientException {
	return ClientException{
		CreateException(message),
	}
}

type NotFoundError struct {
	ClientException
}

func (e NotFoundError) GetErrorCode() int {
	return 404
}

func CreateNotFoundError() NotFoundError {
	return NotFoundError{CreateClientException("Not Found")}
}
