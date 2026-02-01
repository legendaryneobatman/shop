package categories

import (
	"go-shop/pkg/webtool"
	"net/http"
)

type Errors struct {
	BadInput   *webtool.APIError
	CantCreate *webtool.APIError
}

func NewErrors() *Errors {
	return &Errors{
		BadInput:   webtool.NewAPIError("wrong input", http.StatusBadRequest),
		CantCreate: webtool.NewAPIError("cant create category", http.StatusInternalServerError),
	}
}
