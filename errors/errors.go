package errors

type BadRequest struct {
	Code  int
	error error
}

func (e *BadRequest) Error() string {
	return e.error.Error()
}

func NewBadRequestError(err error) *BadRequest {
	return &BadRequest{
		error: err,
		Code:  400,
	}
}

type NotFoundError struct {
	Code  int
	error error
}

func (e *NotFoundError) Error() string {
	return e.error.Error()
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{
		error: err,
		Code:  404,
	}
}

type Error struct {
	Error string `json:"error"`
}
