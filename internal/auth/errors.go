package auth

import (
	"go-shop/pkg/webtool"
	"net/http"
)

type Errors struct {
	UserAlreadyExists  *webtool.APIError
	InvalidCredentials *webtool.APIError
	NoTokenFound       *webtool.APIError
	UserUnauthorized   *webtool.APIError
}

func NewAuthErrors() *Errors {
	return &Errors{
		UserAlreadyExists:  webtool.NewAPIError("user with this username already exists", http.StatusBadRequest),
		InvalidCredentials: webtool.NewAPIError("invalid username or password", http.StatusInternalServerError),
		NoTokenFound:       webtool.NewAPIError("no token found", http.StatusInternalServerError),
		UserUnauthorized:   webtool.NewAPIError("user unauthorized", http.StatusUnauthorized),
	}
}
