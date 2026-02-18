package product

import (
	"shop/pkg/webtool"
	"net/http"
)

var BadInput = webtool.MakeError("wrong input", http.StatusBadRequest)
var NoSlug = webtool.MakeError("empty or no slug provided", http.StatusInternalServerError)
var CantGetProducts = webtool.MakeError("cant get products for category", http.StatusInternalServerError)

