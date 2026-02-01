package webtool

import "github.com/gin-gonic/gin"

func MakeHandler(handlerWithError func(c *gin.Context) *APIError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handlerWithError(c)
		if err != nil && err.Error != nil {
			NewErrorResponse(c, err.HTTPStatus, err.Error.Error())
			return
		}
	}
}
