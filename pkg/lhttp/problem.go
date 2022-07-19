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

// Problem represents an error response from other services
type Problem struct {
	error
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

// New creates a new error with an HTTP status
func New(status int, detail, msg string, args ...interface{}) error {
	return &Problem{
		error:  fmt.Errorf(newErrorFormat, fmt.Sprintf(msg, args...), status),
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	}
}

// Wrap wraps an error with a message and an HTTP status
func Wrap(err error, status int, detail, msg string, args ...interface{}) error {
	return &Problem{
		error:  fmt.Errorf(errorWrapFormat, fmt.Sprintf(msg, args...), status, err),
		Type:   fmt.Sprintf(typeFormat, status),
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	}
}

// ToStatusCode attempts to unwrap the received error to find a status code
func ToStatusCode(err error) int {
	p := &Problem{}
	if ok := errors.As(err, p); !ok {
		return http.StatusInternalServerError
	}

	return p.Status
}

// ToProblem attempts to unwrap the received error to find the wrapped
func ToProblem(err error) *Problem {
	p := &Problem{}
	if ok := errors.As(err, p); !ok {
		p.Type = unknownError
		p.Status = http.StatusInternalServerError
		p.Title = http.StatusText(http.StatusInternalServerError)
	}

	return p
}
