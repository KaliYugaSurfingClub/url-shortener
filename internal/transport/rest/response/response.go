package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"shortener/errs"
)

const (
	statusOk    = "Ok"
	statusError = "errorResponse"
)

type response struct {
	Status string        `json:"status"`
	Data   any           `json:"data,omitempty"`
	Error  errorResponse `json:"error,omitempty"`
}

type errorResponse struct {
	Code       string `json:"code"`
	Validation error  `json:"fields,omitempty"`
}

func Ok(w http.ResponseWriter, data any) {
	resp := response{
		Status: statusOk,
		Data:   data,
	}

	respJson, _ := json.Marshal(resp)
	w.Write(respJson)
}

func Error(w http.ResponseWriter, log *slog.Logger, err error) {
	var e *errs.Error

	switch {
	case err == nil:
		nilErrorResponse(w, log)
	case errors.As(err, &e):
		defaultErrorResponse(w, log, e)
	default:
		unknownErrorResponse(w, log, err)
	}
}

func defaultErrorResponse(w http.ResponseWriter, log *slog.Logger, err *errs.Error) {
	httpStatusCode := httpErrorStatusCode(err.Kind)

	log.Error(
		"error response:",
		slog.Any("stack", errs.OpStack(err)),
		slog.String("msg", err.Error()),
		slog.String("kind", err.Kind.String()),
		slog.String("code", string(err.Code)),
		slog.String("param", string(err.Param)),
		slog.Int("httpCode", httpStatusCode),
	)

	resp := newErrResponse(err)
	errJSON, _ := json.Marshal(resp)

	w.WriteHeader(httpStatusCode)
	w.Write(errJSON)
}

func nilErrorResponse(w http.ResponseWriter, log *slog.Logger) {
	log.Error(
		"Unanticipated nil error - no response body sent",
		slog.Int("HTTP Error StatusCode", http.StatusInternalServerError),
	)

	w.WriteHeader(http.StatusInternalServerError)
}

func unknownErrorResponse(w http.ResponseWriter, log *slog.Logger, err error) {
	resp := response{
		Status: statusError,
		Error: errorResponse{
			Code: errs.Unanticipated.String(),
		},
	}

	log.Error("Unknown Error", slog.String("msg", err.Error()))

	errJSON, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errJSON)
}

func newErrResponse(err *errs.Error) response {
	const internalCode string = "internal server error"
	const validationCode string = "validation error"

	switch err.Kind {
	case errs.Internal, errs.Database:
		return response{
			Status: statusError,
			Error: errorResponse{
				Code: internalCode,
			},
		}
	case errs.Validation:
		return response{
			Status: statusError,
			Error: errorResponse{
				Code:       validationCode,
				Validation: err,
			},
		}
	default:
		code := string(err.Code)
		if code == "" {
			code = err.Kind.String()
		}

		return response{
			Status: statusError,
			Error: errorResponse{
				Code: code,
			},
		}
	}
}

func httpErrorStatusCode(k errs.Kind) int {
	switch k {
	case errs.Invalid, errs.Exist, errs.NotExist, errs.Private, errs.BrokenLink, errs.Validation, errs.InvalidRequest:
		return http.StatusBadRequest
	case errs.UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case errs.Other, errs.IO, errs.Internal, errs.Database, errs.Unanticipated:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
