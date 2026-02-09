package webtool

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code int `json:"code"`
	Source string `json:"source"`
}
type HandlerWithError = func(c *gin.Context) *APIError

func MakeHandler(handlerWithError HandlerWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handlerWithError(c)
		if err != nil && err.Source != nil {
			response := ErrorResponse{
				Message: err.Message,
				Code:    err.Code,
				Source:  err.Source.Error(),
			}

			logrus.WithFields(logrus.Fields{
				"Message": err.Message,
				"Code":    err.Code,
				"Source":  err.Source.Error(),
			}).Errorf("Error in handler")
			c.AbortWithStatusJSON(err.Code, response)
			return
		}
	}
}
