package auth

import (
	"shop/pkg/webtool"
	"net/http"
)


var ErrUserAlreadyExists = func (err error) *webtool.APIError {
	return webtool.NewAPIError("user with this username already exists", http.StatusBadRequest, err)
}
var ErrInvalidCredentials = func (err error) *webtool.APIError {
	return webtool.NewAPIError("invalid username or password", http.StatusBadRequest, err)
}
var ErrNoTokenFound = func (err error) *webtool.APIError {
	return webtool.NewAPIError("no token found", http.StatusInternalServerError, err)
}
var ErrUserUnauthorized = func (err error) *webtool.APIError {
	return webtool.NewAPIError("user unauthorized", http.StatusUnauthorized, err)
}
