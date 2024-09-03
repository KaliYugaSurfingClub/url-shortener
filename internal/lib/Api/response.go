package Api

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
)

//todo refactor

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

// todo pass errors
func Error(msg string) Response {
	return Response{
		Error:  msg,
		Status: StatusError,
	}
}

// todo maybe smthStruct.use("url", field %s is a required field)

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "alphanumunicode":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must contains only letters and nums", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
