package reserr

import (
	"net/http"
)

type HTTPError interface {
	error
	StatusCode() int
	Name() string
	Message() string
}

func NewHTTPError(status int, name string, err error, message string) error {
	return &httpError{status, name, err, message}
}

type httpError struct {
	code int
	name string
	error
	message string
}

func (e *httpError) StatusCode() int {
	return e.code
}

func (e *httpError) Name() string {
	return e.name
}

func (e *httpError) Message() string {
	return e.message
}

type ErrorOutput struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

func Accepted(name string, err error, message string) error {
	return &httpError{http.StatusAccepted, name, err, message}
}

func BadRequest(name string, err error, message string) error {
	return &httpError{http.StatusBadRequest, name, err, message}
}

func Unauthorized(name string, err error, message string) error {
	return &httpError{http.StatusUnauthorized, name, err, message}
}

func NotFound(name string, err error, message string) error {
	return &httpError{http.StatusNotFound, name, err, message}
}

func Conflict(name string, err error, message string) error {
	return &httpError{http.StatusConflict, name, err, message}
}

func Forbidden(name string, err error, message string) error {
	return &httpError{http.StatusForbidden, name, err, message}
}

func Internal(name string, err error, message string) error {
	return &httpError{http.StatusInternalServerError, name, err, message}
}
