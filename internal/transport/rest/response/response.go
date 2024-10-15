package response

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

const (
	ValidationErrorCode = "validation error"
	InternalError       = "internal server error"
)

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

func WithData(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}

func WithOk() Response {
	return Response{
		Status: StatusOk,
	}
}
