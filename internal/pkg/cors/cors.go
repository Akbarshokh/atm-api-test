package cors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, UPDATE, GET, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, multipart/form-data")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Writer.Header().Add("Content-Type", "application/json; charset=utf-8")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
