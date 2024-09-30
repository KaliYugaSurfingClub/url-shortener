package rest

import (
	"errors"
	"shortener/internal/core"
)

type ErrorResponse struct {
	Code string `json:"error code"`
}

func NewErrorResponse(err error) ErrorResponse {
	logicErr := &core.LogicErr{}

	if errors.As(err, logicErr) {
		return ErrorResponse{Code: logicErr.Error()}
	}

	return ErrorResponse{Code: "internal error"}
}
