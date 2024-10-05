package response

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/thoas/go-funk"
)

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

const (
	ValidationErrorCode = "validation error"
	InternalError       = "internal server error"
)

type ValidationError struct {
	Filed string `json:"filed"`
	Msg   string `json:"msg"`
}

func (e ValidationError) Error() string {
	return e.Msg
}

func NewValidationError(filed string, err error) ValidationError {
	return ValidationError{Filed: filed, Msg: err.Error()}
}

type Response struct {
	Status           string            `json:"status"`
	Data             any               `json:"data,omitempty"`
	ErrorCode        string            `json:"error,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

func WithError(err error) (resp Response) {
	return Response{
		Status:    StatusError,
		ErrorCode: err.Error(),
	}
}

func WithInternalError() Response {
	return Response{
		Status:    StatusError,
		ErrorCode: InternalError,
	}
}

func WithValidationErrors(errs error) (resp Response) {
	validationErrors, ok := errs.(validation.Errors)
	if !ok {
		panic("take not validation error")
	}

	return Response{
		Status:           StatusError,
		ErrorCode:        ValidationErrorCode,
		ValidationErrors: funk.Map(validationErrors, NewValidationError).([]ValidationError),
	}
}

func WithOk(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}
