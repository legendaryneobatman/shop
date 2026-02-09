package list

import (
	"go-shop/pkg/webtool"
	"net/http"
)

var NoUserFound = webtool.MakeError("No user found", http.StatusUnauthorized)
var BadRequest = webtool.MakeError("Invalid input", http.StatusBadRequest)
var CantCreate = webtool.MakeError("Failed to create list", http.StatusInternalServerError)
var InvalidPagination = webtool.MakeError("Invalid pagination parameters", http.StatusBadRequest)
var CantFindElements = webtool.MakeError("Failed to find elements", http.StatusInternalServerError)
