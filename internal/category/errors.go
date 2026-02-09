package category

import (
	"go-shop/pkg/webtool"
	"net/http"
)

var BadInput = webtool.MakeError("wrong input", http.StatusBadRequest)
var NoID = webtool.MakeError("No ID provided in input, check query", http.StatusInternalServerError)
var CantCreate = webtool.MakeError("cant create category", http.StatusInternalServerError)
var CantFind = webtool.MakeError("cant find category", http.StatusInternalServerError)
var CantUpdate = webtool.MakeError("cant update category", http.StatusInternalServerError)
var CantDelete = webtool.MakeError("cant delete category", http.StatusInternalServerError)
var CantUpdateOrder = webtool.MakeError("cant update category order", http.StatusInternalServerError)
var CantUpdateStatus = webtool.MakeError("cant update category status", http.StatusInternalServerError)
var CantGetList = webtool.MakeError("cant get category list", http.StatusInternalServerError)
var CantGetTree = webtool.MakeError("cant get tree", http.StatusInternalServerError)
