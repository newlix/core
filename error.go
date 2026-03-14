package core

import (
	"errors"
	"net/http"
)

// ServerError is a server error.
type ServerError struct {
	Status  int
	Message string
}

// Error implementation.
func (e ServerError) Error() string {
	return e.Message
}

// Error returns a new ServerError with HTTP status code, kind and message.
func Error(status int, message string) error {
	return ServerError{
		Status:  status,
		Message: message,
	}
}

// BadRequest returns a new bad request error.
func BadRequest(message string) error {
	return Error(http.StatusBadRequest, message)
}

// WriteError writes an error.
//
// The message in the response uses the Error()
// implementation.
func WriteError(w http.ResponseWriter, err error) {
	c := 500
	var e ServerError
	if errors.As(err, &e) && e.Status != 0 {
		c = e.Status
	}
	http.Error(w, err.Error(), c)
}
