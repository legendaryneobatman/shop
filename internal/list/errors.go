package list

import (
	"go-shop/pkg/webtool"
	"net/http"
)

type Errors struct {
	NoUserFound       *webtool.APIError
	BadRequest        *webtool.APIError
	CantCreate        *webtool.APIError
	CantFindElements  *webtool.APIError
	InvalidPagination *webtool.APIError
}

func NewErrors() *Errors {
	return &Errors{
		NoUserFound:       webtool.NewAPIError("No user found", http.StatusUnauthorized),
		BadRequest:        webtool.NewAPIError("Invalid input", http.StatusBadRequest),
		CantCreate:        webtool.NewAPIError("Failed to create list", http.StatusInternalServerError),
		InvalidPagination: webtool.NewAPIError("Invalid pagination parameters", http.StatusBadRequest),
		CantFindElements:  webtool.NewAPIError("Failed to find elements", http.StatusInternalServerError),
	}
}
