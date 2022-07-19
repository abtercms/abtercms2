// Package lhttp for wrapping errors with HTTP status codes
package lhttp

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	newErrorFormat  = "%s (status %d)"
	errorWrapFormat = "%s (status %d), err: %w"
	unknownError    = "unknown error"
	typeFormat      = "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/%d"
)

// Problem represents an error response from other services.
type Problem struct {
	error
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

// NewProblem creates a new error with an HTTP status.
func NewProblem(status int, detail, msg string, args ...interface{}) error {
	return &Problem{
		error:  fmt.Errorf(newErrorFormat, fmt.Sprintf(msg, args...), status), // nolint: goerr113
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	}
}

// WrapProblem wraps an error with a message and an HTTP status.
func WrapProblem(err error, status int, detail, msg string, args ...interface{}) error {
	return &Problem{
		error:  fmt.Errorf(errorWrapFormat, fmt.Sprintf(msg, args...), status, err),
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	}
}

// ToProblem attempts to unwrap the received error to find the wrapped.
func ToProblem(err error) *Problem {
	var problem *Problem
	if ok := errors.As(err, problem); !ok {
		problem = &Problem{
			error:  err,
			Type:   unknownError,
			Status: http.StatusInternalServerError,
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
	}

	return problem
}
