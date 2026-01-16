// middleware/size_limit.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestSizeLimitMiddleware(maxSizeMB int) gin.HandlerFunc {
	maxBytes := int64(maxSizeMB * 1024 * 1024)

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
