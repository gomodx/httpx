package httpx

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ErrCtx ResMap

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Context ResMap `json:"context"`

	Err    error `json:"-"`
	Status int   `json:"-"`
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e Error) Unwrap() error { return e.Err }

func (e Error) Error() string {
	return fmt.Sprintf("HTTP Error %d %s %s %v", e.Status, e.Code, e.Message, e.Err)
}

func (e Error) WithErrCtx(ctx ErrCtx) Error {
	e.Context = ResMap(ctx)
	return e
}

func (e Error) WithErr(err error) Error {
	if e.Err != nil {
		e.Err = errors.Wrap(err, e.Err.Error())
	} else {
		e.Err = err
	}

	return e
}

func (e Error) WithMessage(msg string) Error {
	e.Message = msg
	return e
}

var (
	ErrInternalServer = Error{
		Message: "internal server error",
		Code:    "INTERNAL_SERVER_ERROR",
		Status:  http.StatusInternalServerError,
	}

	ErrNotFound = Error{
		Message: "resource not found",
		Code:    "NOT_FOUND_ERROR",
		Status:  http.StatusNotFound,
	}

	ErrUnauthorized = Error{
		Message: "unauthorized",
		Code:    "UNAUTHORIZED",
		Status:  http.StatusUnauthorized,
	}

	ErrBadRequest = Error{
		Message: "bad request",
		Code:    "BAD_REQUEST",
		Status:  http.StatusBadRequest,
	}
)
