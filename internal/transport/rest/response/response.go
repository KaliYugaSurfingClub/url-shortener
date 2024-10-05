package response

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
	Status           string `json:"status"`
	Data             any    `json:"data,omitempty"`
	ErrorCode        string `json:"error,omitempty"`
	ValidationErrors error  `json:"validation_errors,omitempty"`
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
	return Response{
		Status:           StatusError,
		ErrorCode:        ValidationErrorCode,
		ValidationErrors: errs,
	}
}

func WithOk(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}
