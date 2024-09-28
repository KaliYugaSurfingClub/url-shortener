package rest

type ErrorResponse struct {
	Code string `json:"error code"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Code: err.Error()}
}
