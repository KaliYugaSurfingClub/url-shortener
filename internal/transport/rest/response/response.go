package response

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type Response struct {
	Status    string `json:"status"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error,omitempty"`
}

func NewError(err error) (resp Response) {
	return Response{
		Status:    StatusError,
		ErrorCode: err.Error(),
	}
}

func NewInternalError() Response {
	return Response{
		Status:    StatusError,
		ErrorCode: "internal server error",
	}
}

func NewOk(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}
