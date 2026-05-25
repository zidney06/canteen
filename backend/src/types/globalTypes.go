package types

type ServiceError struct {
	Message    string
	HttpStatus int
}

func (e *ServiceError) Error() string {
	return e.Message
}

// mungkin seharusnya buat satu type error yang bisa di mengerti oleh semua layer

type ErrorType int

const (
	DbError ErrorType = iota
	NotFound
	ClientError
	ParsingError
	Unauthoried
)

type RepoError struct {
	Message string
	Type    ErrorType
}

func (e *RepoError) Error() string {
	return e.Message
}

func (e *RepoError) ToServiceError() *ServiceError {
	switch e.Type {
	case 1:
		return &ServiceError{
			Message:    e.Message,
			HttpStatus: 404,
		}
	case 2:
		return &ServiceError{
			Message:    e.Message,
			HttpStatus: 400,
		}
	default:
		return &ServiceError{
			Message:    "Internal server error",
			HttpStatus: 500,
		}
	}
}
