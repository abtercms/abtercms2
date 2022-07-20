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
	typeFormat      = "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/%d"
)

var (
	errUnknown = errors.New("unknown error")
)

// Problem represents an error response from other services.
// https://datatracker.ietf.org/doc/html/rfc7807
type Problem struct {
	error
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// NewProblem creates a new error with an HTTP status.
func NewProblem(status int, msg string, args ...interface{}) *Problem {
	return &Problem{
		error:  fmt.Errorf(newErrorFormat, fmt.Sprintf(msg, args...), status), // nolint: goerr113
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: "",
	}
}

// WrapProblem wraps an error with a message and an HTTP status.
func WrapProblem(err error, status int, msg string, args ...interface{}) *Problem {
	return &Problem{
		error:  fmt.Errorf(errorWrapFormat, fmt.Sprintf(msg, args...), status, err),
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: "",
	}
}

// ToProblem attempts to unwrap the received error to find the wrapped.
func ToProblem(err error) *Problem {
	if err == nil {
		err = errUnknown
	}

	problem, ok := err.(*Problem) // nolint: errorlint
	if !ok {
		problem = unwrapProblem(err)
	}

	if problem == nil {
		problem = &Problem{
			error:  fmt.Errorf(errorWrapFormat, err.Error(), http.StatusInternalServerError, err),
			Type:   fmt.Sprintf(typeFormat, http.StatusInternalServerError),
			Status: http.StatusInternalServerError,
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: "",
		}
	}

	if problem.Detail == "" {
		problem.Detail = err.Error()
	}

	return problem
}

func unwrapProblem(err error) *Problem {
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return nil
	}

	problem, ok := unwrapped.(*Problem) // nolint: errorlint
	if !ok {
		return unwrapProblem(unwrapped)
	}

	return problem
}
