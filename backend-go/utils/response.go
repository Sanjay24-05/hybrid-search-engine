// utils/response.go
package utils

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func RespondError(c *gin.Context, status int, err error, message string) {
	c.JSON(status, ErrorResponse{
		Error:   err.Error(),
		Message: message,
	})
}

func RespondSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
