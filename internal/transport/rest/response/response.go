package response

import (
	"errors"
	"shortener/internal/core"
)

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func NewError(err error) Response {
	resp := Response{
		Status: StatusError,
	}

	logicErr := &core.LogicErr{}

	if errors.As(err, logicErr) {
		resp.Error = logicErr.Error()
	} else {
		resp.Error = "internal server error"
	}

	return resp
}

func NewOk(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}
