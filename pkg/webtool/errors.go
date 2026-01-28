package webtool

import "errors"

type APIError struct {
	Error      error
	HTTPStatus int
}

func NewAPIError(message string, status int) *APIError {
	return &APIError{
		Error:      errors.New(message),
		HTTPStatus: status,
	}
}
